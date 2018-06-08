package statefile

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/states"
	"github.com/hashicorp/terraform/tfdiags"
)

func upgradeStateV3ToV4(old *stateV3) (*stateV4, error) {

	if old.Serial < 0 {
		// The new format is using uint64 here, which should be fine for any
		// real state (we only used positive integers in practice) but we'll
		// catch this explicitly here to avoid weird behavior if a state file
		// has been tampered with in some way.
		return nil, fmt.Errorf("state has serial less than zero, which is invalid")
	}

	new := &stateV4{
		TerraformVersion: old.TFVersion,
		Serial:           uint64(old.Serial),
		Lineage:          old.Lineage,
		RootOutputs:      map[string]outputStateV4{},
		Resources:        []resourceStateV4{},
	}

	if new.TerraformVersion == "" {
		// Older formats considered this to be optional, but now it's required
		// and so we'll stub it out with something that's definitely older
		// than the version that really created this state.
		new.TerraformVersion = "0.0.0"
	}

	for _, msOld := range old.Modules {
		if len(msOld.Path) < 1 || msOld.Path[0] != "root" {
			return nil, fmt.Errorf("state contains invalid module path %#v", msOld.Path)
		}

		// Convert legacy-style module address into our newer address type.
		// Since these old formats are only generated by versions of Terraform
		// that don't support count and for_each on modules, we can just assume
		// all of the modules are unkeyed.
		moduleAddr := make(addrs.ModuleInstance, len(msOld.Path)-1)
		for i, name := range msOld.Path[1:] {
			moduleAddr[i] = addrs.ModuleInstanceStep{
				Name:        name,
				InstanceKey: addrs.NoKey,
			}
		}

		// In a v3 state file, a "resource state" is actually an instance
		// state, so we need to fill in a missing level of heirarchy here
		// by lazily creating resource states as we encounter them.
		// We'll track them in here, keyed on the string representation of
		// the resource address.
		resourceStates := map[string]*resourceStateV4{}

		for legacyAddr, rsOld := range msOld.Resources {
			instAddr, err := parseLegacyResourceAddress(legacyAddr)
			if err != nil {
				return nil, err
			}

			resAddr := instAddr.Resource
			rs, exists := resourceStates[resAddr.String()]
			if !exists {
				var modeStr string
				switch resAddr.Mode {
				case addrs.ManagedResourceMode:
					modeStr = "managed"
				case addrs.DataResourceMode:
					modeStr = "data"
				default:
					return nil, fmt.Errorf("state contains resource %s with an unsupported resource mode", resAddr)
				}

				// In state versions prior to 4 we allowed each instance of a
				// resource to have its own provider configuration address,
				// which makes no real sense in practice because providers
				// are associated with resources in the configuration. We
				// elevate that to the resource level during this upgrade,
				// implicitly taking the provider address of the first instance
				// we encounter for each resource. While this is lossy in
				// theory, in practice there is no reason for these values to
				// differ between instances.
				var providerAddr addrs.AbsProviderConfig
				oldProviderAddr := rsOld.Provider
				if strings.Contains(oldProviderAddr, "provider.") {
					// Smells like a new-style provider address, but we'll test it.
					var diags tfdiags.Diagnostics
					providerAddr, diags = addrs.ParseAbsProviderConfigStr(oldProviderAddr)
					if diags.HasErrors() {
						return nil, diags.Err()
					}
				} else {
					// Smells like an old-style module-local provider address,
					// which we'll need to migrate. We'll assume it's referring
					// to the same module the resource is in, which might be
					// incorrect but it'll get fixed up next time any updates
					// are made to an instance.
					if oldProviderAddr != "" {
						localAddr, diags := addrs.ParseProviderConfigCompactStr(oldProviderAddr)
						if diags.HasErrors() {
							return nil, diags.Err()
						}
						providerAddr = localAddr.Absolute(moduleAddr)
					} else {
						providerAddr = resAddr.DefaultProviderConfig().Absolute(moduleAddr)
					}
				}

				rs = &resourceStateV4{
					Module:         moduleAddr.String(),
					Mode:           modeStr,
					Type:           resAddr.Type,
					Name:           resAddr.Name,
					Instances:      []instanceObjectStateV4{},
					ProviderConfig: providerAddr.String(),
				}
				resourceStates[resAddr.String()] = rs
			}

			// Now we'll deal with the instance itself, which may either be
			// the first instance in a resource we just created or an additional
			// instance for a resource added on a prior loop.
			instKey := instAddr.Key
			if isOld := rsOld.Primary; isOld != nil {
				isNew, err := upgradeInstanceObjectV3ToV4(rsOld, isOld, instKey, states.NotDeposed)
				if err != nil {
					return nil, fmt.Errorf("failed to migrate primary generation of %s: %s", instAddr, err)
				}
				rs.Instances = append(rs.Instances, *isNew)
			}
			for i, isOld := range rsOld.Deposed {
				// When we migrate old instances we'll use sequential deposed
				// keys just so that the upgrade result is deterministic. New
				// deposed keys allocated moving forward will be pseudorandomly
				// selected, but we check for collisions and so these
				// non-random ones won't hurt.
				deposedKey := states.DeposedKey(fmt.Sprintf("%08x", i+1))
				isNew, err := upgradeInstanceObjectV3ToV4(rsOld, isOld, instKey, deposedKey)
				if err != nil {
					return nil, fmt.Errorf("failed to migrate deposed generation index %d of %s: %s", i, instAddr, err)
				}
				rs.Instances = append(rs.Instances, *isNew)
			}

			if instKey != addrs.NoKey && rs.EachMode == "" {
				rs.EachMode = "list"
			}
		}

		for _, rs := range resourceStates {
			new.Resources = append(new.Resources, *rs)
		}

		if len(msOld.Path) == 1 && msOld.Path[0] == "root" {
			// We'll migrate the outputs for this module too, then.
			for name, oldOS := range msOld.Outputs {
				newOS := outputStateV4{
					Sensitive: oldOS.Sensitive,
				}

				valRaw := oldOS.Value
				valSrc, err := json.Marshal(valRaw)
				if err != nil {
					// Should never happen, because this value came from JSON
					// in the first place and so we're just round-tripping here.
					return nil, fmt.Errorf("failed to serialize output %q value as JSON: %s", name, err)
				}

				// The "type" field in state V2 wasn't really that useful
				// since it was only able to capture string vs. list vs. map.
				// For this reason, during upgrade we'll just discard it
				// altogether and use cty's idea of the implied type of
				// turning our old value into JSON.
				ty, err := ctyjson.ImpliedType(valSrc)
				if err != nil {
					// REALLY should never happen, because we literally just
					// encoded this as JSON above!
					return nil, fmt.Errorf("failed to parse output %q value from JSON: %s", name, err)
				}

				// ImpliedType tends to produce structural types, but since older
				// version of Terraform didn't support those a collection type
				// is probably what was intended, so we'll see if we can
				// interpret our value as one.
				ty = simplifyImpliedValueType(ty)

				tySrc, err := ctyjson.MarshalType(ty)
				if err != nil {
					return nil, fmt.Errorf("failed to serialize output %q type as JSON: %s", name, err)
				}

				newOS.ValueRaw = json.RawMessage(valSrc)
				newOS.ValueTypeRaw = json.RawMessage(tySrc)

				new.RootOutputs[name] = newOS
			}
		}
	}

	new.normalize()

	return new, nil
}

func upgradeInstanceObjectV3ToV4(rsOld *resourceStateV2, isOld *instanceStateV2, instKey addrs.InstanceKey, deposedKey states.DeposedKey) (*instanceObjectStateV4, error) {

	// Schema versions were, in prior formats, a private concern of the provider
	// SDK, and not a first-class concept in the state format. Here we're
	// sniffing for the pre-0.12 SDK's way of representing schema versions
	// and promoting it to our first-class field if we find it. We'll ignore
	// if if it doesn't look like what the SDK would've written. If this
	// sniffing fails then we'll assume schema version 0.
	var schemaVersion uint64
	migratedSchemaVersion := false
	if raw, exists := isOld.Meta["schema_version"]; exists {
		if rawStr, ok := raw.(string); ok {
			v, err := strconv.ParseUint(rawStr, 10, 64)
			if err == nil {
				schemaVersion = v
				migratedSchemaVersion = true
			}
		}
	}

	private := map[string]interface{}{}
	for k, v := range isOld.Meta {
		if k == "schema_version" && migratedSchemaVersion {
			// We're gonna promote this into our first-class schema version field
			continue
		}
		private[k] = v
	}
	var privateJSON []byte
	if len(private) != 0 {
		var err error
		privateJSON, err = json.Marshal(private)
		if err != nil {
			// This shouldn't happen, because the Meta values all came from JSON
			// originally anyway.
			return nil, fmt.Errorf("cannot serialize private instance object data: %s", err)
		}
	}

	var status string
	if isOld.Tainted {
		status = "tainted"
	}

	var instKeyRaw interface{}
	switch tk := instKey.(type) {
	case addrs.IntKey:
		instKeyRaw = int(tk)
	case addrs.StringKey:
		instKeyRaw = string(tk)
	default:
		if instKeyRaw != nil {
			return nil, fmt.Errorf("insupported instance key: %#v", instKey)
		}
	}

	return &instanceObjectStateV4{
		IndexKey:       instKeyRaw,
		Status:         status,
		Deposed:        string(deposedKey),
		AttributesFlat: isOld.Attributes,
		Dependencies:   rsOld.Dependencies,
		SchemaVersion:  schemaVersion,
		PrivateRaw:     privateJSON,
	}, nil
}

// parseLegacyResourceAddress parses the different identifier format used
// state formats before version 4, like "instance.name.0".
func parseLegacyResourceAddress(s string) (addrs.ResourceInstance, error) {
	var ret addrs.ResourceInstance

	// Split based on ".". Every resource address should have at least two
	// elements (type and name).
	parts := strings.Split(s, ".")
	if len(parts) < 2 || len(parts) > 4 {
		return ret, fmt.Errorf("invalid internal resource address format: %s", s)
	}

	// Data resource if we have at least 3 parts and the first one is data
	ret.Resource.Mode = addrs.ManagedResourceMode
	if len(parts) > 2 && parts[0] == "data" {
		ret.Resource.Mode = addrs.DataResourceMode
		parts = parts[1:]
	}

	// If we're not a data resource and we have more than 3, then it is an error
	if len(parts) > 3 && ret.Resource.Mode != addrs.DataResourceMode {
		return ret, fmt.Errorf("invalid internal resource address format: %s", s)
	}

	// Build the parts of the resource address that are guaranteed to exist
	ret.Resource.Type = parts[0]
	ret.Resource.Name = parts[1]
	ret.Key = addrs.NoKey

	// If we have more parts, then we have an index. Parse that.
	if len(parts) > 2 {
		idx, err := strconv.ParseInt(parts[2], 0, 0)
		if err != nil {
			return ret, fmt.Errorf("error parsing resource address %q: %s", s, err)
		}

		ret.Key = addrs.IntKey(idx)
	}

	return ret, nil
}

// simplifyImpliedValueType attempts to heuristically simplify a value type
// derived from a legacy stored output value into something simpler that
// is closer to what would've fitted into the pre-v0.12 value type system.
func simplifyImpliedValueType(ty cty.Type) cty.Type {
	switch {
	case ty.IsTupleType():
		// If all of the element types are the same then we'll make this
		// a list instead. This is very likely to be true, since prior versions
		// of Terraform did not officially support mixed-type collections.

		if ty.Equals(cty.EmptyTuple) {
			// Don't know what the element type would be, then.
			return ty
		}

		etys := ty.TupleElementTypes()
		ety := etys[0]
		for _, other := range etys[1:] {
			if !other.Equals(ety) {
				// inconsistent types
				return ty
			}
		}
		ety = simplifyImpliedValueType(ety)
		return cty.List(ety)

	case ty.IsObjectType():
		// If all of the attribute types are the same then we'll make this
		// a map instead. This is very likely to be true, since prior versions
		// of Terraform did not officially support mixed-type collections.

		if ty.Equals(cty.EmptyObject) {
			// Don't know what the element type would be, then.
			return ty
		}

		atys := ty.AttributeTypes()
		var ety cty.Type
		for _, other := range atys {
			if ety == cty.NilType {
				ety = other
				continue
			}
			if !other.Equals(ety) {
				// inconsistent types
				return ty
			}
		}
		ety = simplifyImpliedValueType(ety)
		return cty.Map(ety)

	default:
		// No other normalizations are possible
		return ty
	}
}
