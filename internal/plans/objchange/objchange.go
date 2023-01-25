package objchange

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/internal/configs/configschema"
)

// ProposedNew constructs a proposed new object value by combining the
// computed attribute values from "prior" with the configured attribute values
// from "config".
//
// Both value must conform to the given schema's implied type, or this function
// will panic.
//
// The prior value must be wholly known, but the config value may be unknown
// or have nested unknown values.
//
// The merging of the two objects includes the attributes of any nested blocks,
// which will be correlated in a manner appropriate for their nesting mode.
// Note in particular that the correlation for blocks backed by sets is a
// heuristic based on matching non-computed attribute values and so it may
// produce strange results with more "extreme" cases, such as a nested set
// block where _all_ attributes are computed.
func ProposedNew(schema *configschema.Block, prior, config cty.Value) cty.Value {
	// If the config and prior are both null, return early here before
	// populating the prior block. The prevents non-null blocks from appearing
	// the proposed state value.
	if config.IsNull() && prior.IsNull() {
		return prior
	}

	if prior.IsNull() {
		// In this case, we will construct a synthetic prior value that is
		// similar to the result of decoding an empty configuration block,
		// which simplifies our handling of the top-level attributes/blocks
		// below by giving us one non-null level of object to pull values from.
		//
		// "All attributes null" happens to be the definition of EmptyValue for
		// a Block, so we can just delegate to that
		prior = schema.EmptyValue()
	}
	return proposedNew(schema, prior, config)
}

// PlannedDataResourceObject is similar to proposedNewBlock but tailored for
// planning data resources in particular. Specifically, it replaces the values
// of any Computed attributes not set in the configuration with an unknown
// value, which serves as a placeholder for a value to be filled in by the
// provider when the data resource is finally read.
//
// Data resources are different because the planning of them is handled
// entirely within Terraform Core and not subject to customization by the
// provider. This function is, in effect, producing an equivalent result to
// passing the proposedNewBlock result into a provider's PlanResourceChange
// function, assuming a fixed implementation of PlanResourceChange that just
// fills in unknown values as needed.
func PlannedDataResourceObject(schema *configschema.Block, config cty.Value) cty.Value {
	// Our trick here is to run the proposedNewBlock logic with an
	// entirely-unknown prior value. Because of cty's unknown short-circuit
	// behavior, any operation on prior returns another unknown, and so
	// unknown values propagate into all of the parts of the resulting value
	// that would normally be filled in by preserving the prior state.
	prior := cty.UnknownVal(schema.ImpliedType())
	return proposedNew(schema, prior, config)
}

func proposedNew(schema *configschema.Block, prior, config cty.Value) cty.Value {
	if config.IsNull() || !config.IsKnown() {
		// A block config should never be null at this point. The only nullable
		// block type is NestingSingle, which will return early before coming
		// back here. We'll allow the null here anyway to free callers from
		// needing to specifically check for these cases, and any mismatch will
		// be caught in validation, so just take the prior value rather than
		// the invalid null.
		return prior
	}

	if (!prior.Type().IsObjectType()) || (!config.Type().IsObjectType()) {
		panic("ProposedNew only supports object-typed values")
	}

	// From this point onwards, we can assume that both values are non-null
	// object types, and that the config value itself is known (though it
	// may contain nested values that are unknown.)
	newAttrs := proposedNewAttributes(schema.Attributes, prior, config)

	// Merging nested blocks is a little more complex, since we need to
	// correlate blocks between both objects and then recursively propose
	// a new object for each. The correlation logic depends on the nesting
	// mode for each block type.
	for name, blockType := range schema.BlockTypes {
		priorV := prior.GetAttr(name)
		configV := config.GetAttr(name)
		newAttrs[name] = proposedNewNestedBlock(blockType, priorV, configV)
	}

	return cty.ObjectVal(newAttrs)
}

// proposedNewBlockOrObject dispatched the schema to either ProposedNew or
// proposedNewObjectAttributes depending on the given type.
func proposedNewBlockOrObject(schema nestedSchema, prior, config cty.Value) cty.Value {
	switch schema := schema.(type) {
	case *configschema.Block:
		return ProposedNew(schema, prior, config)
	case *configschema.Object:
		return proposedNewObjectAttributes(schema, prior, config)
	default:
		panic(fmt.Sprintf("unexpected schema type %T", schema))
	}
}

func proposedNewNestedBlock(schema *configschema.NestedBlock, prior, config cty.Value) cty.Value {
	// The only time we should encounter an entirely unknown block is from the
	// use of dynamic with an unknown for_each expression.
	if !config.IsKnown() {
		return config
	}

	newV := config

	switch schema.Nesting {
	case configschema.NestingSingle:
		// A NestingSingle configuration block value can be null, and since it
		// cannot be computed we can always take the configuration value.
		if config.IsNull() {
			break
		}

		// Otherwise use the same assignment rules as NestingGroup
		fallthrough
	case configschema.NestingGroup:
		newV = ProposedNew(&schema.Block, prior, config)

	case configschema.NestingList:
		newV = proposedNewNestingList(&schema.Block, prior, config)

	case configschema.NestingMap:
		newV = proposedNewNestingMap(&schema.Block, prior, config)

	case configschema.NestingSet:
		newV = proposedNewNestingSet(&schema.Block, prior, config)

	default:
		// Should never happen, since the above cases are comprehensive.
		panic(fmt.Sprintf("unsupported block nesting mode %s", schema.Nesting))
	}

	return newV
}

func proposedNewNestedType(schema *configschema.Object, prior, config cty.Value) cty.Value {
	// if the config isn't known at all, then we must use that value
	if !config.IsKnown() {
		return config
	}

	// Even if the config is null or empty, we will be using this default value.
	newV := config

	switch schema.Nesting {
	case configschema.NestingSingle:
		// If the config is null, we already have our value. If the attribute
		// is optional+computed, we won't reach this branch with a null value
		// since the computed case would have been taken.
		if config.IsNull() {
			break
		}

		newV = proposedNewObjectAttributes(schema, prior, config)

	case configschema.NestingList:
		newV = proposedNewNestingList(schema, prior, config)

	case configschema.NestingMap:
		newV = proposedNewNestingMap(schema, prior, config)

	case configschema.NestingSet:
		newV = proposedNewNestingSet(schema, prior, config)

	default:
		// Should never happen, since the above cases are comprehensive.
		panic(fmt.Sprintf("unsupported attribute nesting mode %s", schema.Nesting))
	}

	return newV
}

func proposedNewNestingList(schema nestedSchema, prior, config cty.Value) cty.Value {
	newV := config

	// Nested blocks are correlated by index.
	configVLen := 0
	if !config.IsNull() {
		configVLen = config.LengthInt()
	}
	if configVLen > 0 {
		newVals := make([]cty.Value, 0, configVLen)
		for it := config.ElementIterator(); it.Next(); {
			idx, configEV := it.Element()
			if prior.IsKnown() && (prior.IsNull() || !prior.HasIndex(idx).True()) {
				// If there is no corresponding prior element then
				// we just take the config value as-is.
				newVals = append(newVals, configEV)
				continue
			}
			priorEV := prior.Index(idx)

			newVals = append(newVals, proposedNewBlockOrObject(schema, priorEV, configEV))
		}
		// Despite the name, a NestingList might also be a tuple, if
		// its nested schema contains dynamically-typed attributes.
		if config.Type().IsTupleType() {
			newV = cty.TupleVal(newVals)
		} else {
			newV = cty.ListVal(newVals)
		}
	}

	return newV
}

func proposedNewNestingMap(schema nestedSchema, prior, config cty.Value) cty.Value {
	newV := config

	newVals := map[string]cty.Value{}

	if config.IsNull() || !config.IsKnown() || config.LengthInt() == 0 {
		// We already assigned newVal and there's nothing to compare in
		// config.
		return newV
	}
	cfgMap := config.AsValueMap()

	// prior may be null or empty
	priorMap := map[string]cty.Value{}
	if !prior.IsNull() && prior.IsKnown() && prior.LengthInt() > 0 {
		priorMap = prior.AsValueMap()
	}

	for name, configEV := range cfgMap {
		priorEV, inPrior := priorMap[name]
		if !inPrior {
			// If there is no corresponding prior element then
			// we just take the config value as-is.
			newVals[name] = configEV
			continue
		}

		newVals[name] = proposedNewBlockOrObject(schema, priorEV, configEV)
	}

	// The value must leave as the same type it came in as
	switch {
	case config.Type().IsObjectType():
		// Although we call the nesting mode "map", we actually use
		// object values so that elements might have different types
		// in case of dynamically-typed attributes.
		newV = cty.ObjectVal(newVals)
	default:
		newV = cty.MapVal(newVals)
	}

	return newV
}

func proposedNewNestingSet(schema nestedSchema, prior, config cty.Value) cty.Value {
	newV := config

	if !config.Type().IsSetType() {
		panic("configschema.NestingSet value is not a set as expected")
	}

	// Nested set elements are correlated by comparing the element values after
	// eliminating all of the computed attributes. In practice, this means that
	// any config change produces an entirely new nested object, and we only
	// propagate prior computed values if the non-computed attribute values are
	// identical.
	var cmpVals [][2]cty.Value
	if prior.IsKnown() && !prior.IsNull() {
		cmpVals = setElementCompareValues(schema, prior)
	}

	configVLen := 0
	if config.IsKnown() && !config.IsNull() {
		configVLen = config.LengthInt()
	}
	if configVLen > 0 {
		used := make([]bool, len(cmpVals)) // track used elements in case multiple have the same compare value
		newVals := make([]cty.Value, 0, configVLen)
		for it := config.ElementIterator(); it.Next(); {
			_, configEV := it.Element()
			var priorEV cty.Value
			for i, cmp := range cmpVals {
				if used[i] {
					continue
				}

				if cmp[1].RawEquals(configEV) {
					priorEV = cmp[0]
					used[i] = true // we can't use this value on a future iteration
					break
				}
			}
			if priorEV == cty.NilVal {
				priorEV = cty.NullVal(config.Type().ElementType())
			}

			newVals = append(newVals, proposedNewBlockOrObject(schema, priorEV, configEV))
		}
		newV = cty.SetVal(newVals)
	}

	return newV
}

func proposedNewObjectAttributes(schema *configschema.Object, prior, config cty.Value) cty.Value {
	if config.IsNull() {
		return config
	}

	return cty.ObjectVal(proposedNewAttributes(schema.Attributes, prior, config))
}

func proposedNewAttributes(attrs map[string]*configschema.Attribute, prior, config cty.Value) map[string]cty.Value {
	newAttrs := make(map[string]cty.Value, len(attrs))
	for name, attr := range attrs {
		var priorV cty.Value
		if prior.IsNull() {
			priorV = cty.NullVal(prior.Type().AttributeType(name))
		} else {
			priorV = prior.GetAttr(name)
		}

		configV := config.GetAttr(name)

		var newV cty.Value
		switch {
		// required isn't considered when constructing the plan, so attributes
		// are essentially either computed or not computed. In the case of
		// optional+computed, they are only computed when there is no
		// configuration.
		case attr.Computed && configV.IsNull():
			// configV will always be null in this case, by definition.
			// priorV may also be null, but that's okay.
			newV = priorV

			// the exception to the above is that if the config is optional and
			// the _prior_ value contains non-computed values, we can infer
			// that the config must have been non-null previously.
			if optionalValueNotComputable(attr, priorV) {
				newV = configV
			}

		case attr.NestedType != nil:
			// For non-computed NestedType attributes, we need to descend
			// into the individual nested attributes to build the final
			// value, unless the entire nested attribute is unknown.
			newV = proposedNewNestedType(attr.NestedType, priorV, configV)
		default:
			// For non-computed attributes, we always take the config value,
			// even if it is null. If it's _required_ then null values
			// should've been caught during an earlier validation step, and
			// so we don't really care about that here.
			newV = configV
		}
		newAttrs[name] = newV
	}
	return newAttrs
}

// setElementCompareValues takes a known, non-null value of a cty.Set type and
// returns a table -- constructed of two-element arrays -- that maps original
// set element values to corresponding values that have all of the computed
// values removed, making them suitable for comparison with values obtained
// from configuration. The element type of the set must conform to the implied
// type of the given schema, or this function will panic.
//
// In the resulting slice, the zeroth element of each array is the original
// value and the one-indexed element is the corresponding "compare value".
//
// This is intended to help correlate prior elements with configured elements
// in proposedNewBlock. The result is a heuristic rather than an exact science,
// since e.g. two separate elements may reduce to the same value through this
// process. The caller must therefore be ready to deal with duplicates.
func setElementCompareValues(schema nestedSchema, set cty.Value) [][2]cty.Value {
	ret := make([][2]cty.Value, 0, set.LengthInt())
	for it := set.ElementIterator(); it.Next(); {
		_, ev := it.Element()
		ret = append(ret, [2]cty.Value{ev, setElementComputedAsNull(schema, ev)})
	}
	return ret
}

// nestedSchema is used as a generic container for either a
// *configschema.Object, or *configschema.Block.
type nestedSchema interface {
	AttributeByPath(cty.Path) *configschema.Attribute
}

func setElementComputedAsNull(schema nestedSchema, elem cty.Value) cty.Value {
	elem, _ = cty.Transform(elem, func(path cty.Path, v cty.Value) (cty.Value, error) {
		if v.IsNull() || !v.IsKnown() {
			return v, nil
		}

		attr := schema.AttributeByPath(path)
		if attr == nil {
			// we must be at a nested block rather than an attribute, continue
			return v, nil
		}

		if attr.Computed {
			// If this could have been computed, we always return null in order
			// to compare config and state values.
			//
			// The only way to determine if an optional+computed value is
			// optional or computed is see if it is set in the config, however
			// for set elements we don't yet know which element value
			// corresponds to the configuration.
			return cty.NullVal(v.Type()), nil
		}
		return v, nil
	})

	return elem
}

// optionalValueNotComputable is used to check if an object in state must
// have at least partially come from configuration. If the prior value has any
// non-null attributes which are not computed in the schema, then we know there
// was previously a configuration value which set those.
//
// This is used when the configuration contains a null optional+computed value,
// and we want to know if we should plan to send the null value or the prior
// state.
func optionalValueNotComputable(schema *configschema.Attribute, val cty.Value) bool {
	if !schema.Optional {
		return false
	}

	// We must have a NestedType for complex nested attributes in order
	// to find nested computed values in the first place.
	if schema.NestedType == nil {
		return false
	}

	foundNonComputedAttr := false
	cty.Walk(val, func(path cty.Path, v cty.Value) (bool, error) {
		if v.IsNull() {
			return true, nil
		}

		attr := schema.NestedType.AttributeByPath(path)
		if attr == nil {
			return true, nil
		}

		if !attr.Computed {
			foundNonComputedAttr = true
			return false, nil
		}
		return true, nil
	})

	return foundNonComputedAttr
}
