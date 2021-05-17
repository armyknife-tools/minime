package planfile

import (
	"fmt"
	"io"
	"io/ioutil"

	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/internal/tfdiags"
	"github.com/hashicorp/terraform/plans"
	"github.com/hashicorp/terraform/plans/internal/planproto"
	"github.com/hashicorp/terraform/states"
	"github.com/hashicorp/terraform/version"
	"github.com/zclconf/go-cty/cty"
)

const tfplanFormatVersion = 3
const tfplanFilename = "tfplan"

// ---------------------------------------------------------------------------
// This file deals with the internal structure of the "tfplan" sub-file within
// the plan file format. It's all private API, wrapped by methods defined
// elsewhere. This is the only file that should import the
// ../internal/planproto package, which contains the ugly stubs generated
// by the protobuf compiler.
// ---------------------------------------------------------------------------

// readTfplan reads a protobuf-encoded description from the plan portion of
// a plan file, which is stored in a special file in the archive called
// "tfplan".
func readTfplan(r io.Reader) (*plans.Plan, error) {
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var rawPlan planproto.Plan
	err = proto.Unmarshal(src, &rawPlan)
	if err != nil {
		return nil, fmt.Errorf("parse error: %s", err)
	}

	if rawPlan.Version != tfplanFormatVersion {
		return nil, fmt.Errorf("unsupported plan file format version %d; only version %d is supported", rawPlan.Version, tfplanFormatVersion)
	}

	if rawPlan.TerraformVersion != version.String() {
		return nil, fmt.Errorf("plan file was created by Terraform %s, but this is %s; plan files cannot be transferred between different Terraform versions", rawPlan.TerraformVersion, version.String())
	}

	plan := &plans.Plan{
		VariableValues: map[string]plans.DynamicValue{},
		Changes: &plans.Changes{
			Outputs:   []*plans.OutputChangeSrc{},
			Resources: []*plans.ResourceInstanceChangeSrc{},
		},

		ProviderSHA256s: map[string][]byte{},
	}

	switch rawPlan.UiMode {
	case planproto.Mode_NORMAL:
		plan.UIMode = plans.NormalMode
	case planproto.Mode_DESTROY:
		plan.UIMode = plans.DestroyMode
	case planproto.Mode_REFRESH_ONLY:
		plan.UIMode = plans.RefreshOnlyMode
	default:
		return nil, fmt.Errorf("plan has invalid mode %s", rawPlan.UiMode)
	}

	for _, rawOC := range rawPlan.OutputChanges {
		name := rawOC.Name
		change, err := changeFromTfplan(rawOC.Change)
		if err != nil {
			return nil, fmt.Errorf("invalid plan for output %q: %s", name, err)
		}

		plan.Changes.Outputs = append(plan.Changes.Outputs, &plans.OutputChangeSrc{
			// All output values saved in the plan file are root module outputs,
			// since we don't retain others. (They can be easily recomputed
			// during apply).
			Addr:      addrs.OutputValue{Name: name}.Absolute(addrs.RootModuleInstance),
			ChangeSrc: *change,
			Sensitive: rawOC.Sensitive,
		})
	}

	for _, rawRC := range rawPlan.ResourceChanges {
		change, err := resourceChangeFromTfplan(rawRC)
		if err != nil {
			// errors from resourceChangeFromTfplan already include context
			return nil, err
		}

		plan.Changes.Resources = append(plan.Changes.Resources, change)
	}

	for _, rawTargetAddr := range rawPlan.TargetAddrs {
		target, diags := addrs.ParseTargetStr(rawTargetAddr)
		if diags.HasErrors() {
			return nil, fmt.Errorf("plan contains invalid target address %q: %s", target, diags.Err())
		}
		plan.TargetAddrs = append(plan.TargetAddrs, target.Subject)
	}

	for _, rawReplaceAddr := range rawPlan.ForceReplaceAddrs {
		addr, diags := addrs.ParseAbsResourceInstanceStr(rawReplaceAddr)
		if diags.HasErrors() {
			return nil, fmt.Errorf("plan contains invalid force-replace address %q: %s", addr, diags.Err())
		}
		plan.ForceReplaceAddrs = append(plan.ForceReplaceAddrs, addr)
	}

	for name, rawHashObj := range rawPlan.ProviderHashes {
		if len(rawHashObj.Sha256) == 0 {
			return nil, fmt.Errorf("no SHA256 hash for provider %q plugin", name)
		}

		plan.ProviderSHA256s[name] = rawHashObj.Sha256
	}

	for name, rawVal := range rawPlan.Variables {
		val, err := valueFromTfplan(rawVal)
		if err != nil {
			return nil, fmt.Errorf("invalid value for input variable %q: %s", name, err)
		}
		plan.VariableValues[name] = val
	}

	if rawBackend := rawPlan.Backend; rawBackend == nil {
		return nil, fmt.Errorf("plan file has no backend settings; backend settings are required")
	} else {
		config, err := valueFromTfplan(rawBackend.Config)
		if err != nil {
			return nil, fmt.Errorf("plan file has invalid backend configuration: %s", err)
		}
		plan.Backend = plans.Backend{
			Type:      rawBackend.Type,
			Config:    config,
			Workspace: rawBackend.Workspace,
		}
	}

	return plan, nil
}

func resourceChangeFromTfplan(rawChange *planproto.ResourceInstanceChange) (*plans.ResourceInstanceChangeSrc, error) {
	if rawChange == nil {
		// Should never happen in practice, since protobuf can't represent
		// a nil value in a list.
		return nil, fmt.Errorf("resource change object is absent")
	}

	ret := &plans.ResourceInstanceChangeSrc{}

	moduleAddr := addrs.RootModuleInstance
	if rawChange.ModulePath != "" {
		var diags tfdiags.Diagnostics
		moduleAddr, diags = addrs.ParseModuleInstanceStr(rawChange.ModulePath)
		if diags.HasErrors() {
			return nil, diags.Err()
		}
	}

	providerAddr, diags := addrs.ParseAbsProviderConfigStr(rawChange.Provider)
	if diags.HasErrors() {
		return nil, diags.Err()
	}
	ret.ProviderAddr = providerAddr

	var mode addrs.ResourceMode
	switch rawChange.Mode {
	case planproto.ResourceInstanceChange_managed:
		mode = addrs.ManagedResourceMode
	case planproto.ResourceInstanceChange_data:
		mode = addrs.DataResourceMode
	default:
		return nil, fmt.Errorf("resource has invalid mode %s", rawChange.Mode)
	}

	typeName := rawChange.Type
	name := rawChange.Name

	resAddr := addrs.Resource{
		Mode: mode,
		Type: typeName,
		Name: name,
	}

	var instKey addrs.InstanceKey
	switch rawTk := rawChange.InstanceKey.(type) {
	case nil:
	case *planproto.ResourceInstanceChange_Int:
		instKey = addrs.IntKey(rawTk.Int)
	case *planproto.ResourceInstanceChange_Str:
		instKey = addrs.StringKey(rawTk.Str)
	default:
		return nil, fmt.Errorf("instance of %s has invalid key type %T", resAddr.Absolute(moduleAddr), rawChange.InstanceKey)
	}

	ret.Addr = resAddr.Instance(instKey).Absolute(moduleAddr)

	if rawChange.DeposedKey != "" {
		if len(rawChange.DeposedKey) != 8 {
			return nil, fmt.Errorf("deposed object for %s has invalid deposed key %q", ret.Addr, rawChange.DeposedKey)
		}
		ret.DeposedKey = states.DeposedKey(rawChange.DeposedKey)
	}

	ret.RequiredReplace = cty.NewPathSet()
	for _, p := range rawChange.RequiredReplace {
		path, err := pathFromTfplan(p)
		if err != nil {
			return nil, fmt.Errorf("invalid path in required replace: %s", err)
		}
		ret.RequiredReplace.Add(path)
	}

	change, err := changeFromTfplan(rawChange.Change)
	if err != nil {
		return nil, fmt.Errorf("invalid plan for resource %s: %s", ret.Addr, err)
	}

	ret.ChangeSrc = *change

	switch rawChange.ActionReason {
	case planproto.ResourceInstanceActionReason_NONE:
		ret.ActionReason = plans.ResourceInstanceChangeNoReason
	case planproto.ResourceInstanceActionReason_REPLACE_BECAUSE_CANNOT_UPDATE:
		ret.ActionReason = plans.ResourceInstanceReplaceBecauseCannotUpdate
	case planproto.ResourceInstanceActionReason_REPLACE_BECAUSE_TAINTED:
		ret.ActionReason = plans.ResourceInstanceReplaceBecauseTainted
	case planproto.ResourceInstanceActionReason_REPLACE_BY_REQUEST:
		ret.ActionReason = plans.ResourceInstanceReplaceByRequest
	default:
		return nil, fmt.Errorf("resource has invalid action reason %s", rawChange.ActionReason)
	}

	if len(rawChange.Private) != 0 {
		ret.Private = rawChange.Private
	}

	return ret, nil
}

func changeFromTfplan(rawChange *planproto.Change) (*plans.ChangeSrc, error) {
	if rawChange == nil {
		return nil, fmt.Errorf("change object is absent")
	}

	ret := &plans.ChangeSrc{}

	// -1 indicates that there is no index. We'll customize these below
	// depending on the change action, and then decode.
	beforeIdx, afterIdx := -1, -1

	switch rawChange.Action {
	case planproto.Action_NOOP:
		ret.Action = plans.NoOp
		beforeIdx = 0
		afterIdx = 0
	case planproto.Action_CREATE:
		ret.Action = plans.Create
		afterIdx = 0
	case planproto.Action_READ:
		ret.Action = plans.Read
		beforeIdx = 0
		afterIdx = 1
	case planproto.Action_UPDATE:
		ret.Action = plans.Update
		beforeIdx = 0
		afterIdx = 1
	case planproto.Action_DELETE:
		ret.Action = plans.Delete
		beforeIdx = 0
	case planproto.Action_CREATE_THEN_DELETE:
		ret.Action = plans.CreateThenDelete
		beforeIdx = 0
		afterIdx = 1
	case planproto.Action_DELETE_THEN_CREATE:
		ret.Action = plans.DeleteThenCreate
		beforeIdx = 0
		afterIdx = 1
	default:
		return nil, fmt.Errorf("invalid change action %s", rawChange.Action)
	}

	if beforeIdx != -1 {
		if l := len(rawChange.Values); l <= beforeIdx {
			return nil, fmt.Errorf("incorrect number of values (%d) for %s change", l, rawChange.Action)
		}
		var err error
		ret.Before, err = valueFromTfplan(rawChange.Values[beforeIdx])
		if err != nil {
			return nil, fmt.Errorf("invalid \"before\" value: %s", err)
		}
		if ret.Before == nil {
			return nil, fmt.Errorf("missing \"before\" value: %s", err)
		}
	}
	if afterIdx != -1 {
		if l := len(rawChange.Values); l <= afterIdx {
			return nil, fmt.Errorf("incorrect number of values (%d) for %s change", l, rawChange.Action)
		}
		var err error
		ret.After, err = valueFromTfplan(rawChange.Values[afterIdx])
		if err != nil {
			return nil, fmt.Errorf("invalid \"after\" value: %s", err)
		}
		if ret.After == nil {
			return nil, fmt.Errorf("missing \"after\" value: %s", err)
		}
	}

	sensitive := cty.NewValueMarks("sensitive")
	beforeValMarks, err := pathValueMarksFromTfplan(rawChange.BeforeSensitivePaths, sensitive)
	if err != nil {
		return nil, fmt.Errorf("failed to decode before sensitive paths: %s", err)
	}
	afterValMarks, err := pathValueMarksFromTfplan(rawChange.AfterSensitivePaths, sensitive)
	if err != nil {
		return nil, fmt.Errorf("failed to decode after sensitive paths: %s", err)
	}
	if len(beforeValMarks) > 0 {
		ret.BeforeValMarks = beforeValMarks
	}
	if len(afterValMarks) > 0 {
		ret.AfterValMarks = afterValMarks
	}

	return ret, nil
}

func valueFromTfplan(rawV *planproto.DynamicValue) (plans.DynamicValue, error) {
	if len(rawV.Msgpack) == 0 { // len(0) because that's the default value for a "bytes" in protobuf
		return nil, fmt.Errorf("dynamic value does not have msgpack serialization")
	}

	return plans.DynamicValue(rawV.Msgpack), nil
}

// writeTfplan serializes the given plan into the protobuf-based format used
// for the "tfplan" portion of a plan file.
func writeTfplan(plan *plans.Plan, w io.Writer) error {
	if plan == nil {
		return fmt.Errorf("cannot write plan file for nil plan")
	}
	if plan.Changes == nil {
		return fmt.Errorf("cannot write plan file with nil changeset")
	}

	rawPlan := &planproto.Plan{
		Version:          tfplanFormatVersion,
		TerraformVersion: version.String(),
		ProviderHashes:   map[string]*planproto.Hash{},

		Variables:       map[string]*planproto.DynamicValue{},
		OutputChanges:   []*planproto.OutputChange{},
		ResourceChanges: []*planproto.ResourceInstanceChange{},
	}

	switch plan.UIMode {
	case plans.NormalMode:
		rawPlan.UiMode = planproto.Mode_NORMAL
	case plans.DestroyMode:
		rawPlan.UiMode = planproto.Mode_DESTROY
	case plans.RefreshOnlyMode:
		rawPlan.UiMode = planproto.Mode_REFRESH_ONLY
	default:
		return fmt.Errorf("plan has unsupported mode %s", plan.UIMode)
	}

	for _, oc := range plan.Changes.Outputs {
		// When serializing a plan we only retain the root outputs, since
		// changes to these are externally-visible side effects (e.g. via
		// terraform_remote_state).
		if !oc.Addr.Module.IsRoot() {
			continue
		}

		name := oc.Addr.OutputValue.Name

		// Writing outputs as cty.DynamicPseudoType forces the stored values
		// to also contain dynamic type information, so we can recover the
		// original type when we read the values back in readTFPlan.
		protoChange, err := changeToTfplan(&oc.ChangeSrc)
		if err != nil {
			return fmt.Errorf("cannot write output value %q: %s", name, err)
		}

		rawPlan.OutputChanges = append(rawPlan.OutputChanges, &planproto.OutputChange{
			Name:      name,
			Change:    protoChange,
			Sensitive: oc.Sensitive,
		})
	}

	for _, rc := range plan.Changes.Resources {
		rawRC, err := resourceChangeToTfplan(rc)
		if err != nil {
			return err
		}
		rawPlan.ResourceChanges = append(rawPlan.ResourceChanges, rawRC)
	}

	for _, targetAddr := range plan.TargetAddrs {
		rawPlan.TargetAddrs = append(rawPlan.TargetAddrs, targetAddr.String())
	}

	for _, replaceAddr := range plan.ForceReplaceAddrs {
		rawPlan.ForceReplaceAddrs = append(rawPlan.ForceReplaceAddrs, replaceAddr.String())
	}

	for name, hash := range plan.ProviderSHA256s {
		rawPlan.ProviderHashes[name] = &planproto.Hash{
			Sha256: hash,
		}
	}

	for name, val := range plan.VariableValues {
		rawPlan.Variables[name] = valueToTfplan(val)
	}

	if plan.Backend.Type == "" || plan.Backend.Config == nil {
		// This suggests a bug in the code that created the plan, since it
		// ought to always have a backend populated, even if it's the default
		// "local" backend with a local state file.
		return fmt.Errorf("plan does not have a backend configuration")
	}

	rawPlan.Backend = &planproto.Backend{
		Type:      plan.Backend.Type,
		Config:    valueToTfplan(plan.Backend.Config),
		Workspace: plan.Backend.Workspace,
	}

	src, err := proto.Marshal(rawPlan)
	if err != nil {
		return fmt.Errorf("serialization error: %s", err)
	}

	_, err = w.Write(src)
	if err != nil {
		return fmt.Errorf("failed to write plan to plan file: %s", err)
	}

	return nil
}

func resourceChangeToTfplan(change *plans.ResourceInstanceChangeSrc) (*planproto.ResourceInstanceChange, error) {
	ret := &planproto.ResourceInstanceChange{}

	ret.ModulePath = change.Addr.Module.String()

	relAddr := change.Addr.Resource

	switch relAddr.Resource.Mode {
	case addrs.ManagedResourceMode:
		ret.Mode = planproto.ResourceInstanceChange_managed
	case addrs.DataResourceMode:
		ret.Mode = planproto.ResourceInstanceChange_data
	default:
		return nil, fmt.Errorf("resource %s has unsupported mode %s", relAddr, relAddr.Resource.Mode)
	}

	ret.Type = relAddr.Resource.Type
	ret.Name = relAddr.Resource.Name

	switch tk := relAddr.Key.(type) {
	case nil:
		// Nothing to do, then.
	case addrs.IntKey:
		ret.InstanceKey = &planproto.ResourceInstanceChange_Int{
			Int: int64(tk),
		}
	case addrs.StringKey:
		ret.InstanceKey = &planproto.ResourceInstanceChange_Str{
			Str: string(tk),
		}
	default:
		return nil, fmt.Errorf("resource %s has unsupported instance key type %T", relAddr, relAddr.Key)
	}

	ret.DeposedKey = string(change.DeposedKey)
	ret.Provider = change.ProviderAddr.String()

	requiredReplace := change.RequiredReplace.List()
	ret.RequiredReplace = make([]*planproto.Path, 0, len(requiredReplace))
	for _, p := range requiredReplace {
		path, err := pathToTfplan(p)
		if err != nil {
			return nil, fmt.Errorf("invalid path in required replace: %s", err)
		}
		ret.RequiredReplace = append(ret.RequiredReplace, path)
	}

	valChange, err := changeToTfplan(&change.ChangeSrc)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize resource %s change: %s", relAddr, err)
	}
	ret.Change = valChange

	switch change.ActionReason {
	case plans.ResourceInstanceChangeNoReason:
		ret.ActionReason = planproto.ResourceInstanceActionReason_NONE
	case plans.ResourceInstanceReplaceBecauseCannotUpdate:
		ret.ActionReason = planproto.ResourceInstanceActionReason_REPLACE_BECAUSE_CANNOT_UPDATE
	case plans.ResourceInstanceReplaceBecauseTainted:
		ret.ActionReason = planproto.ResourceInstanceActionReason_REPLACE_BECAUSE_TAINTED
	case plans.ResourceInstanceReplaceByRequest:
		ret.ActionReason = planproto.ResourceInstanceActionReason_REPLACE_BY_REQUEST
	default:
		return nil, fmt.Errorf("resource %s has unsupported action reason %s", relAddr, change.ActionReason)
	}

	if len(change.Private) > 0 {
		ret.Private = change.Private
	}

	return ret, nil
}

func changeToTfplan(change *plans.ChangeSrc) (*planproto.Change, error) {
	ret := &planproto.Change{}

	before := valueToTfplan(change.Before)
	after := valueToTfplan(change.After)

	beforeSensitivePaths, err := pathValueMarksToTfplan(change.BeforeValMarks)
	if err != nil {
		return nil, err
	}
	afterSensitivePaths, err := pathValueMarksToTfplan(change.AfterValMarks)
	if err != nil {
		return nil, err
	}
	ret.BeforeSensitivePaths = beforeSensitivePaths
	ret.AfterSensitivePaths = afterSensitivePaths

	switch change.Action {
	case plans.NoOp:
		ret.Action = planproto.Action_NOOP
		ret.Values = []*planproto.DynamicValue{before} // before and after should be identical
	case plans.Create:
		ret.Action = planproto.Action_CREATE
		ret.Values = []*planproto.DynamicValue{after}
	case plans.Read:
		ret.Action = planproto.Action_READ
		ret.Values = []*planproto.DynamicValue{before, after}
	case plans.Update:
		ret.Action = planproto.Action_UPDATE
		ret.Values = []*planproto.DynamicValue{before, after}
	case plans.Delete:
		ret.Action = planproto.Action_DELETE
		ret.Values = []*planproto.DynamicValue{before}
	case plans.DeleteThenCreate:
		ret.Action = planproto.Action_DELETE_THEN_CREATE
		ret.Values = []*planproto.DynamicValue{before, after}
	case plans.CreateThenDelete:
		ret.Action = planproto.Action_CREATE_THEN_DELETE
		ret.Values = []*planproto.DynamicValue{before, after}
	default:
		return nil, fmt.Errorf("invalid change action %s", change.Action)
	}

	return ret, nil
}

func valueToTfplan(val plans.DynamicValue) *planproto.DynamicValue {
	if val == nil {
		// protobuf can't represent nil, so we'll represent it as a
		// DynamicValue that has no serializations at all.
		return &planproto.DynamicValue{}
	}
	return &planproto.DynamicValue{
		Msgpack: []byte(val),
	}
}

func pathValueMarksFromTfplan(paths []*planproto.Path, marks cty.ValueMarks) ([]cty.PathValueMarks, error) {
	ret := make([]cty.PathValueMarks, 0, len(paths))
	for _, p := range paths {
		path, err := pathFromTfplan(p)
		if err != nil {
			return nil, err
		}
		ret = append(ret, cty.PathValueMarks{
			Path:  path,
			Marks: marks,
		})
	}
	return ret, nil
}

func pathValueMarksToTfplan(pvm []cty.PathValueMarks) ([]*planproto.Path, error) {
	ret := make([]*planproto.Path, 0, len(pvm))
	for _, p := range pvm {
		path, err := pathToTfplan(p.Path)
		if err != nil {
			return nil, err
		}
		ret = append(ret, path)
	}
	return ret, nil
}

func pathFromTfplan(path *planproto.Path) (cty.Path, error) {
	ret := make([]cty.PathStep, 0, len(path.Steps))
	for _, step := range path.Steps {
		switch s := step.Selector.(type) {
		case *planproto.Path_Step_ElementKey:
			dynamicVal, err := valueFromTfplan(s.ElementKey)
			if err != nil {
				return nil, fmt.Errorf("error decoding path index step: %s", err)
			}
			ty, err := dynamicVal.ImpliedType()
			if err != nil {
				return nil, fmt.Errorf("error determining path index type: %s", err)
			}
			val, err := dynamicVal.Decode(ty)
			if err != nil {
				return nil, fmt.Errorf("error decoding path index value: %s", err)
			}
			ret = append(ret, cty.IndexStep{Key: val})
		case *planproto.Path_Step_AttributeName:
			ret = append(ret, cty.GetAttrStep{Name: s.AttributeName})
		default:
			return nil, fmt.Errorf("Unsupported path step %t", step.Selector)
		}
	}
	return ret, nil
}

func pathToTfplan(path cty.Path) (*planproto.Path, error) {
	steps := make([]*planproto.Path_Step, 0, len(path))
	for _, step := range path {
		switch s := step.(type) {
		case cty.IndexStep:
			value, err := plans.NewDynamicValue(s.Key, s.Key.Type())
			if err != nil {
				return nil, fmt.Errorf("Error encoding path step: %s", err)
			}
			steps = append(steps, &planproto.Path_Step{
				Selector: &planproto.Path_Step_ElementKey{
					ElementKey: valueToTfplan(value),
				},
			})
		case cty.GetAttrStep:
			steps = append(steps, &planproto.Path_Step{
				Selector: &planproto.Path_Step_AttributeName{
					AttributeName: s.Name,
				},
			})
		default:
			return nil, fmt.Errorf("Unsupported path step %#v (%t)", step, step)
		}
	}
	return &planproto.Path{
		Steps: steps,
	}, nil
}
