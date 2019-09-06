// Code generated by protoc-gen-go. DO NOT EDIT.
// source: planfile.proto

package planproto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Action describes the type of action planned for an object.
// Not all action values are valid for all object types.
type Action int32

const (
	Action_NOOP               Action = 0
	Action_CREATE             Action = 1
	Action_READ               Action = 2
	Action_UPDATE             Action = 3
	Action_DELETE             Action = 5
	Action_DELETE_THEN_CREATE Action = 6
	Action_CREATE_THEN_DELETE Action = 7
)

var Action_name = map[int32]string{
	0: "NOOP",
	1: "CREATE",
	2: "READ",
	3: "UPDATE",
	5: "DELETE",
	6: "DELETE_THEN_CREATE",
	7: "CREATE_THEN_DELETE",
}

var Action_value = map[string]int32{
	"NOOP":               0,
	"CREATE":             1,
	"READ":               2,
	"UPDATE":             3,
	"DELETE":             5,
	"DELETE_THEN_CREATE": 6,
	"CREATE_THEN_DELETE": 7,
}

func (x Action) String() string {
	return proto.EnumName(Action_name, int32(x))
}

func (Action) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_02431083a6706c5b, []int{0}
}

type ResourceInstanceChange_ResourceMode int32

const (
	ResourceInstanceChange_managed ResourceInstanceChange_ResourceMode = 0
	ResourceInstanceChange_data    ResourceInstanceChange_ResourceMode = 1
)

var ResourceInstanceChange_ResourceMode_name = map[int32]string{
	0: "managed",
	1: "data",
}

var ResourceInstanceChange_ResourceMode_value = map[string]int32{
	"managed": 0,
	"data":    1,
}

func (x ResourceInstanceChange_ResourceMode) String() string {
	return proto.EnumName(ResourceInstanceChange_ResourceMode_name, int32(x))
}

func (ResourceInstanceChange_ResourceMode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_02431083a6706c5b, []int{3, 0}
}

// Plan is the root message type for the tfplan file
type Plan struct {
	// Version is incremented whenever there is a breaking change to
	// the serialization format. Programs reading serialized plans should
	// verify that version is set to the expected value and abort processing
	// if not. A breaking change is any change that may cause an older
	// consumer to interpret the structure incorrectly. This number will
	// not be incremented if an existing consumer can either safely ignore
	// changes to the format or if an existing consumer would fail to process
	// the file for another message- or field-specific reason.
	Version uint64 `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	// The variables that were set when creating the plan. Each value is
	// a msgpack serialization of an HCL value.
	Variables map[string]*DynamicValue `protobuf:"bytes,2,rep,name=variables,proto3" json:"variables,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// An unordered set of proposed changes to resources throughout the
	// configuration, including any nested modules. Use the address of
	// each resource to determine which module it belongs to.
	ResourceChanges []*ResourceInstanceChange `protobuf:"bytes,3,rep,name=resource_changes,json=resourceChanges,proto3" json:"resource_changes,omitempty"`
	// An unordered set of proposed changes to outputs in the root module
	// of the configuration. This set also includes "no action" changes for
	// outputs that are not changing, as context for detecting inconsistencies
	// at apply time.
	OutputChanges []*OutputChange `protobuf:"bytes,4,rep,name=output_changes,json=outputChanges,proto3" json:"output_changes,omitempty"`
	// An unordered set of target addresses to include when applying. If no
	// target addresses are present, the plan applies to the whole
	// configuration.
	TargetAddrs []string `protobuf:"bytes,5,rep,name=target_addrs,json=targetAddrs,proto3" json:"target_addrs,omitempty"`
	// The version string for the Terraform binary that created this plan.
	TerraformVersion string `protobuf:"bytes,14,opt,name=terraform_version,json=terraformVersion,proto3" json:"terraform_version,omitempty"`
	// SHA256 digests of all of the provider plugin binaries that were used
	// in the creation of this plan.
	ProviderHashes map[string]*Hash `protobuf:"bytes,15,rep,name=provider_hashes,json=providerHashes,proto3" json:"provider_hashes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Backend is a description of the backend configuration and other related
	// settings at the time the plan was created.
	Backend              *Backend `protobuf:"bytes,13,opt,name=backend,proto3" json:"backend,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Plan) Reset()         { *m = Plan{} }
func (m *Plan) String() string { return proto.CompactTextString(m) }
func (*Plan) ProtoMessage()    {}
func (*Plan) Descriptor() ([]byte, []int) {
	return fileDescriptor_02431083a6706c5b, []int{0}
}

func (m *Plan) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Plan.Unmarshal(m, b)
}
func (m *Plan) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Plan.Marshal(b, m, deterministic)
}
func (m *Plan) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Plan.Merge(m, src)
}
func (m *Plan) XXX_Size() int {
	return xxx_messageInfo_Plan.Size(m)
}
func (m *Plan) XXX_DiscardUnknown() {
	xxx_messageInfo_Plan.DiscardUnknown(m)
}

var xxx_messageInfo_Plan proto.InternalMessageInfo

func (m *Plan) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *Plan) GetVariables() map[string]*DynamicValue {
	if m != nil {
		return m.Variables
	}
	return nil
}

func (m *Plan) GetResourceChanges() []*ResourceInstanceChange {
	if m != nil {
		return m.ResourceChanges
	}
	return nil
}

func (m *Plan) GetOutputChanges() []*OutputChange {
	if m != nil {
		return m.OutputChanges
	}
	return nil
}

func (m *Plan) GetTargetAddrs() []string {
	if m != nil {
		return m.TargetAddrs
	}
	return nil
}

func (m *Plan) GetTerraformVersion() string {
	if m != nil {
		return m.TerraformVersion
	}
	return ""
}

func (m *Plan) GetProviderHashes() map[string]*Hash {
	if m != nil {
		return m.ProviderHashes
	}
	return nil
}

func (m *Plan) GetBackend() *Backend {
	if m != nil {
		return m.Backend
	}
	return nil
}

// Backend is a description of backend configuration and other related settings.
type Backend struct {
	Type                 string        `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Config               *DynamicValue `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	Workspace            string        `protobuf:"bytes,3,opt,name=workspace,proto3" json:"workspace,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Backend) Reset()         { *m = Backend{} }
func (m *Backend) String() string { return proto.CompactTextString(m) }
func (*Backend) ProtoMessage()    {}
func (*Backend) Descriptor() ([]byte, []int) {
	return fileDescriptor_02431083a6706c5b, []int{1}
}

func (m *Backend) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Backend.Unmarshal(m, b)
}
func (m *Backend) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Backend.Marshal(b, m, deterministic)
}
func (m *Backend) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backend.Merge(m, src)
}
func (m *Backend) XXX_Size() int {
	return xxx_messageInfo_Backend.Size(m)
}
func (m *Backend) XXX_DiscardUnknown() {
	xxx_messageInfo_Backend.DiscardUnknown(m)
}

var xxx_messageInfo_Backend proto.InternalMessageInfo

func (m *Backend) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Backend) GetConfig() *DynamicValue {
	if m != nil {
		return m.Config
	}
	return nil
}

func (m *Backend) GetWorkspace() string {
	if m != nil {
		return m.Workspace
	}
	return ""
}

// Change represents a change made to some object, transforming it from an old
// state to a new state.
type Change struct {
	// Not all action values are valid for all object types. Consult
	// the documentation for any message that embeds Change.
	Action Action `protobuf:"varint,1,opt,name=action,proto3,enum=tfplan.Action" json:"action,omitempty"`
	// msgpack-encoded HCL values involved in the change.
	// - For update and replace, two values are provided that give the old and new values,
	//   respectively.
	// - For create, one value is provided that gives the new value to be created
	// - For delete, one value is provided that describes the value being deleted
	// - For read, two values are provided that give the prior value for this object
	//   (or null, if no prior value exists) and the value that was or will be read,
	//   respectively.
	// - For no-op, one value is provided that is left unmodified by this non-change.
	Values               []*DynamicValue `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Change) Reset()         { *m = Change{} }
func (m *Change) String() string { return proto.CompactTextString(m) }
func (*Change) ProtoMessage()    {}
func (*Change) Descriptor() ([]byte, []int) {
	return fileDescriptor_02431083a6706c5b, []int{2}
}

func (m *Change) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Change.Unmarshal(m, b)
}
func (m *Change) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Change.Marshal(b, m, deterministic)
}
func (m *Change) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Change.Merge(m, src)
}
func (m *Change) XXX_Size() int {
	return xxx_messageInfo_Change.Size(m)
}
func (m *Change) XXX_DiscardUnknown() {
	xxx_messageInfo_Change.DiscardUnknown(m)
}

var xxx_messageInfo_Change proto.InternalMessageInfo

func (m *Change) GetAction() Action {
	if m != nil {
		return m.Action
	}
	return Action_NOOP
}

func (m *Change) GetValues() []*DynamicValue {
	if m != nil {
		return m.Values
	}
	return nil
}

type ResourceInstanceChange struct {
	// module_path is an address to the module that defined this resource.
	// module_path is omitted for resources in the root module. For descendent modules
	// it is a string like module.foo.module.bar as would be seen at the beginning of a
	// resource address. The format of this string is not yet frozen and so external
	// callers should treat it as an opaque key for filtering purposes.
	ModulePath string `protobuf:"bytes,1,opt,name=module_path,json=modulePath,proto3" json:"module_path,omitempty"`
	// mode is the resource mode.
	Mode ResourceInstanceChange_ResourceMode `protobuf:"varint,2,opt,name=mode,proto3,enum=tfplan.ResourceInstanceChange_ResourceMode" json:"mode,omitempty"`
	// type is the resource type name, like "aws_instance".
	Type string `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	// name is the logical name of the resource as defined in configuration.
	// For example, in aws_instance.foo this would be "foo".
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	// instance_key is either an integer index or a string key, depending on which iteration
	// attributes ("count" or "for_each") are being used for this resource. If none
	// are in use, this field is omitted.
	//
	// Types that are valid to be assigned to InstanceKey:
	//	*ResourceInstanceChange_Str
	//	*ResourceInstanceChange_Int
	InstanceKey isResourceInstanceChange_InstanceKey `protobuf_oneof:"instance_key"`
	// deposed_key, if set, indicates that this change applies to a deposed
	// object for the indicated instance with the given deposed key. If not
	// set, the change applies to the instance's current object.
	DeposedKey string `protobuf:"bytes,7,opt,name=deposed_key,json=deposedKey,proto3" json:"deposed_key,omitempty"`
	// provider is the address of the provider configuration that this change
	// was planned with, and thus the configuration that must be used to
	// apply it.
	Provider string `protobuf:"bytes,8,opt,name=provider,proto3" json:"provider,omitempty"`
	// Description of the proposed change. May use "create", "read", "update",
	// "replace" and "delete" actions. "no-op" changes are not currently used here
	// but consumers must accept and discard them to allow for future expansion.
	Change *Change `protobuf:"bytes,9,opt,name=change,proto3" json:"change,omitempty"`
	// raw blob value provided by the provider as additional context for the
	// change. Must be considered an opaque value for any consumer other than
	// the provider that generated it, and will be returned verbatim to the
	// provider during the subsequent apply operation.
	Private []byte `protobuf:"bytes,10,opt,name=private,proto3" json:"private,omitempty"`
	// An unordered set of paths that prompted the change action to be
	// "replace" rather than "update". Empty for any action other than
	// "replace".
	RequiredReplace      []*Path  `protobuf:"bytes,11,rep,name=required_replace,json=requiredReplace,proto3" json:"required_replace,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResourceInstanceChange) Reset()         { *m = ResourceInstanceChange{} }
func (m *ResourceInstanceChange) String() string { return proto.CompactTextString(m) }
func (*ResourceInstanceChange) ProtoMessage()    {}
func (*ResourceInstanceChange) Descriptor() ([]byte, []int) {
	return fileDescriptor_02431083a6706c5b, []int{3}
}

func (m *ResourceInstanceChange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResourceInstanceChange.Unmarshal(m, b)
}
func (m *ResourceInstanceChange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResourceInstanceChange.Marshal(b, m, deterministic)
}
func (m *ResourceInstanceChange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResourceInstanceChange.Merge(m, src)
}
func (m *ResourceInstanceChange) XXX_Size() int {
	return xxx_messageInfo_ResourceInstanceChange.Size(m)
}
func (m *ResourceInstanceChange) XXX_DiscardUnknown() {
	xxx_messageInfo_ResourceInstanceChange.DiscardUnknown(m)
}

var xxx_messageInfo_ResourceInstanceChange proto.InternalMessageInfo

func (m *ResourceInstanceChange) GetModulePath() string {
	if m != nil {
		return m.ModulePath
	}
	return ""
}

func (m *ResourceInstanceChange) GetMode() ResourceInstanceChange_ResourceMode {
	if m != nil {
		return m.Mode
	}
	return ResourceInstanceChange_managed
}

func (m *ResourceInstanceChange) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *ResourceInstanceChange) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type isResourceInstanceChange_InstanceKey interface {
	isResourceInstanceChange_InstanceKey()
}

type ResourceInstanceChange_Str struct {
	Str string `protobuf:"bytes,5,opt,name=str,proto3,oneof"`
}

type ResourceInstanceChange_Int struct {
	Int int64 `protobuf:"varint,6,opt,name=int,proto3,oneof"`
}

func (*ResourceInstanceChange_Str) isResourceInstanceChange_InstanceKey() {}

func (*ResourceInstanceChange_Int) isResourceInstanceChange_InstanceKey() {}

func (m *ResourceInstanceChange) GetInstanceKey() isResourceInstanceChange_InstanceKey {
	if m != nil {
		return m.InstanceKey
	}
	return nil
}

func (m *ResourceInstanceChange) GetStr() string {
	if x, ok := m.GetInstanceKey().(*ResourceInstanceChange_Str); ok {
		return x.Str
	}
	return ""
}

func (m *ResourceInstanceChange) GetInt() int64 {
	if x, ok := m.GetInstanceKey().(*ResourceInstanceChange_Int); ok {
		return x.Int
	}
	return 0
}

func (m *ResourceInstanceChange) GetDeposedKey() string {
	if m != nil {
		return m.DeposedKey
	}
	return ""
}

func (m *ResourceInstanceChange) GetProvider() string {
	if m != nil {
		return m.Provider
	}
	return ""
}

func (m *ResourceInstanceChange) GetChange() *Change {
	if m != nil {
		return m.Change
	}
	return nil
}

func (m *ResourceInstanceChange) GetPrivate() []byte {
	if m != nil {
		return m.Private
	}
	return nil
}

func (m *ResourceInstanceChange) GetRequiredReplace() []*Path {
	if m != nil {
		return m.RequiredReplace
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*ResourceInstanceChange) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*ResourceInstanceChange_Str)(nil),
		(*ResourceInstanceChange_Int)(nil),
	}
}

type OutputChange struct {
	// Name of the output as defined in the root module.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Description of the proposed change. May use "no-op", "create",
	// "update" and "delete" actions.
	Change *Change `protobuf:"bytes,2,opt,name=change,proto3" json:"change,omitempty"`
	// Sensitive, if true, indicates that one or more of the values given
	// in "change" is sensitive and should not be shown directly in any
	// rendered plan.
	Sensitive            bool     `protobuf:"varint,3,opt,name=sensitive,proto3" json:"sensitive,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OutputChange) Reset()         { *m = OutputChange{} }
func (m *OutputChange) String() string { return proto.CompactTextString(m) }
func (*OutputChange) ProtoMessage()    {}
func (*OutputChange) Descriptor() ([]byte, []int) {
	return fileDescriptor_02431083a6706c5b, []int{4}
}

func (m *OutputChange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OutputChange.Unmarshal(m, b)
}
func (m *OutputChange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OutputChange.Marshal(b, m, deterministic)
}
func (m *OutputChange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OutputChange.Merge(m, src)
}
func (m *OutputChange) XXX_Size() int {
	return xxx_messageInfo_OutputChange.Size(m)
}
func (m *OutputChange) XXX_DiscardUnknown() {
	xxx_messageInfo_OutputChange.DiscardUnknown(m)
}

var xxx_messageInfo_OutputChange proto.InternalMessageInfo

func (m *OutputChange) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *OutputChange) GetChange() *Change {
	if m != nil {
		return m.Change
	}
	return nil
}

func (m *OutputChange) GetSensitive() bool {
	if m != nil {
		return m.Sensitive
	}
	return false
}

// DynamicValue represents a value whose type is not decided until runtime,
// often based on schema information obtained from a plugin.
//
// At present dynamic values are always encoded as msgpack, with extension
// id 0 used to represent the special "unknown" value indicating results
// that won't be known until after apply.
//
// In future other serialization formats may be used, possibly with a
// transitional period of including both as separate attributes of this type.
// Consumers must ignore attributes they don't support and fail if no supported
// attribute is present. The top-level format version will not be incremented
// for changes to the set of dynamic serialization formats.
type DynamicValue struct {
	Msgpack              []byte   `protobuf:"bytes,1,opt,name=msgpack,proto3" json:"msgpack,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DynamicValue) Reset()         { *m = DynamicValue{} }
func (m *DynamicValue) String() string { return proto.CompactTextString(m) }
func (*DynamicValue) ProtoMessage()    {}
func (*DynamicValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_02431083a6706c5b, []int{5}
}

func (m *DynamicValue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DynamicValue.Unmarshal(m, b)
}
func (m *DynamicValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DynamicValue.Marshal(b, m, deterministic)
}
func (m *DynamicValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DynamicValue.Merge(m, src)
}
func (m *DynamicValue) XXX_Size() int {
	return xxx_messageInfo_DynamicValue.Size(m)
}
func (m *DynamicValue) XXX_DiscardUnknown() {
	xxx_messageInfo_DynamicValue.DiscardUnknown(m)
}

var xxx_messageInfo_DynamicValue proto.InternalMessageInfo

func (m *DynamicValue) GetMsgpack() []byte {
	if m != nil {
		return m.Msgpack
	}
	return nil
}

// Hash represents a hash value.
//
// At present hashes always use the SHA256 algorithm. In future other hash
// algorithms may be used, possibly with a transitional period of including
// both as separate attributes of this type. Consumers must ignore attributes
// they don't support and fail if no supported attribute is present. The
// top-level format version will not be incremented for changes to the set of
// hash algorithms.
type Hash struct {
	Sha256               []byte   `protobuf:"bytes,1,opt,name=sha256,proto3" json:"sha256,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Hash) Reset()         { *m = Hash{} }
func (m *Hash) String() string { return proto.CompactTextString(m) }
func (*Hash) ProtoMessage()    {}
func (*Hash) Descriptor() ([]byte, []int) {
	return fileDescriptor_02431083a6706c5b, []int{6}
}

func (m *Hash) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Hash.Unmarshal(m, b)
}
func (m *Hash) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Hash.Marshal(b, m, deterministic)
}
func (m *Hash) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Hash.Merge(m, src)
}
func (m *Hash) XXX_Size() int {
	return xxx_messageInfo_Hash.Size(m)
}
func (m *Hash) XXX_DiscardUnknown() {
	xxx_messageInfo_Hash.DiscardUnknown(m)
}

var xxx_messageInfo_Hash proto.InternalMessageInfo

func (m *Hash) GetSha256() []byte {
	if m != nil {
		return m.Sha256
	}
	return nil
}

// Path represents a set of steps to traverse into a data structure. It is
// used to refer to a sub-structure within a dynamic data structure presented
// separately.
type Path struct {
	Steps                []*Path_Step `protobuf:"bytes,1,rep,name=steps,proto3" json:"steps,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Path) Reset()         { *m = Path{} }
func (m *Path) String() string { return proto.CompactTextString(m) }
func (*Path) ProtoMessage()    {}
func (*Path) Descriptor() ([]byte, []int) {
	return fileDescriptor_02431083a6706c5b, []int{7}
}

func (m *Path) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Path.Unmarshal(m, b)
}
func (m *Path) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Path.Marshal(b, m, deterministic)
}
func (m *Path) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Path.Merge(m, src)
}
func (m *Path) XXX_Size() int {
	return xxx_messageInfo_Path.Size(m)
}
func (m *Path) XXX_DiscardUnknown() {
	xxx_messageInfo_Path.DiscardUnknown(m)
}

var xxx_messageInfo_Path proto.InternalMessageInfo

func (m *Path) GetSteps() []*Path_Step {
	if m != nil {
		return m.Steps
	}
	return nil
}

type Path_Step struct {
	// Types that are valid to be assigned to Selector:
	//	*Path_Step_AttributeName
	//	*Path_Step_ElementKey
	Selector             isPath_Step_Selector `protobuf_oneof:"selector"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Path_Step) Reset()         { *m = Path_Step{} }
func (m *Path_Step) String() string { return proto.CompactTextString(m) }
func (*Path_Step) ProtoMessage()    {}
func (*Path_Step) Descriptor() ([]byte, []int) {
	return fileDescriptor_02431083a6706c5b, []int{7, 0}
}

func (m *Path_Step) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Path_Step.Unmarshal(m, b)
}
func (m *Path_Step) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Path_Step.Marshal(b, m, deterministic)
}
func (m *Path_Step) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Path_Step.Merge(m, src)
}
func (m *Path_Step) XXX_Size() int {
	return xxx_messageInfo_Path_Step.Size(m)
}
func (m *Path_Step) XXX_DiscardUnknown() {
	xxx_messageInfo_Path_Step.DiscardUnknown(m)
}

var xxx_messageInfo_Path_Step proto.InternalMessageInfo

type isPath_Step_Selector interface {
	isPath_Step_Selector()
}

type Path_Step_AttributeName struct {
	AttributeName string `protobuf:"bytes,1,opt,name=attribute_name,json=attributeName,proto3,oneof"`
}

type Path_Step_ElementKey struct {
	ElementKey *DynamicValue `protobuf:"bytes,2,opt,name=element_key,json=elementKey,proto3,oneof"`
}

func (*Path_Step_AttributeName) isPath_Step_Selector() {}

func (*Path_Step_ElementKey) isPath_Step_Selector() {}

func (m *Path_Step) GetSelector() isPath_Step_Selector {
	if m != nil {
		return m.Selector
	}
	return nil
}

func (m *Path_Step) GetAttributeName() string {
	if x, ok := m.GetSelector().(*Path_Step_AttributeName); ok {
		return x.AttributeName
	}
	return ""
}

func (m *Path_Step) GetElementKey() *DynamicValue {
	if x, ok := m.GetSelector().(*Path_Step_ElementKey); ok {
		return x.ElementKey
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Path_Step) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Path_Step_AttributeName)(nil),
		(*Path_Step_ElementKey)(nil),
	}
}

func init() {
	proto.RegisterEnum("tfplan.Action", Action_name, Action_value)
	proto.RegisterEnum("tfplan.ResourceInstanceChange_ResourceMode", ResourceInstanceChange_ResourceMode_name, ResourceInstanceChange_ResourceMode_value)
	proto.RegisterType((*Plan)(nil), "tfplan.Plan")
	proto.RegisterMapType((map[string]*Hash)(nil), "tfplan.Plan.ProviderHashesEntry")
	proto.RegisterMapType((map[string]*DynamicValue)(nil), "tfplan.Plan.VariablesEntry")
	proto.RegisterType((*Backend)(nil), "tfplan.Backend")
	proto.RegisterType((*Change)(nil), "tfplan.Change")
	proto.RegisterType((*ResourceInstanceChange)(nil), "tfplan.ResourceInstanceChange")
	proto.RegisterType((*OutputChange)(nil), "tfplan.OutputChange")
	proto.RegisterType((*DynamicValue)(nil), "tfplan.DynamicValue")
	proto.RegisterType((*Hash)(nil), "tfplan.Hash")
	proto.RegisterType((*Path)(nil), "tfplan.Path")
	proto.RegisterType((*Path_Step)(nil), "tfplan.Path.Step")
}

func init() { proto.RegisterFile("planfile.proto", fileDescriptor_02431083a6706c5b) }

var fileDescriptor_02431083a6706c5b = []byte{
	// 893 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x55, 0xe1, 0x6e, 0xe3, 0x44,
	0x10, 0xae, 0x63, 0xc7, 0x49, 0x26, 0xa9, 0x9b, 0x5b, 0x50, 0x65, 0x95, 0xd3, 0x11, 0x2c, 0xc1,
	0x85, 0x3b, 0x94, 0x4a, 0x41, 0x50, 0x0e, 0x7e, 0xa0, 0xf6, 0x1a, 0x29, 0xd5, 0x41, 0x1b, 0x2d,
	0xa5, 0x3f, 0xf8, 0x81, 0xb5, 0xb1, 0xa7, 0x89, 0x55, 0xc7, 0x36, 0xbb, 0x9b, 0xa0, 0x3c, 0x10,
	0x0f, 0xc1, 0x4b, 0xf0, 0x4c, 0x68, 0x77, 0x6d, 0x27, 0x95, 0x7a, 0xfd, 0x95, 0x9d, 0x6f, 0x66,
	0x3e, 0xcf, 0x7e, 0x33, 0xb3, 0x01, 0xaf, 0x48, 0x59, 0x76, 0x9f, 0xa4, 0x38, 0x2a, 0x78, 0x2e,
	0x73, 0xe2, 0xca, 0x7b, 0x85, 0x04, 0xff, 0x39, 0xe0, 0xcc, 0x52, 0x96, 0x11, 0x1f, 0x5a, 0x1b,
	0xe4, 0x22, 0xc9, 0x33, 0xdf, 0x1a, 0x58, 0x43, 0x87, 0x56, 0x26, 0x79, 0x07, 0x9d, 0x0d, 0xe3,
	0x09, 0x9b, 0xa7, 0x28, 0xfc, 0xc6, 0xc0, 0x1e, 0x76, 0xc7, 0x9f, 0x8d, 0x4c, 0xfa, 0x48, 0xa5,
	0x8e, 0xee, 0x2a, 0xef, 0x24, 0x93, 0x7c, 0x4b, 0x77, 0xd1, 0xe4, 0x0a, 0xfa, 0x1c, 0x45, 0xbe,
	0xe6, 0x11, 0x86, 0xd1, 0x92, 0x65, 0x0b, 0x14, 0xbe, 0xad, 0x19, 0x5e, 0x55, 0x0c, 0xb4, 0xf4,
	0x5f, 0x65, 0x42, 0xb2, 0x2c, 0xc2, 0xf7, 0x3a, 0x8c, 0x1e, 0x55, 0x79, 0xc6, 0x16, 0xe4, 0x27,
	0xf0, 0xf2, 0xb5, 0x2c, 0xd6, 0xb2, 0x26, 0x72, 0x34, 0xd1, 0xa7, 0x15, 0xd1, 0x8d, 0xf6, 0x96,
	0xe9, 0x87, 0xf9, 0x9e, 0x25, 0xc8, 0x17, 0xd0, 0x93, 0x8c, 0x2f, 0x50, 0x86, 0x2c, 0x8e, 0xb9,
	0xf0, 0x9b, 0x03, 0x7b, 0xd8, 0xa1, 0x5d, 0x83, 0x9d, 0x2b, 0x88, 0xbc, 0x85, 0x17, 0x12, 0x39,
	0x67, 0xf7, 0x39, 0x5f, 0x85, 0x95, 0x12, 0xde, 0xc0, 0x1a, 0x76, 0x68, 0xbf, 0x76, 0xdc, 0x95,
	0x92, 0x5c, 0xc1, 0x51, 0xc1, 0xf3, 0x4d, 0x12, 0x23, 0x0f, 0x97, 0x4c, 0x2c, 0x51, 0xf8, 0x47,
	0xba, 0x9a, 0xc1, 0x23, 0x61, 0x66, 0x65, 0xcc, 0x54, 0x87, 0x18, 0x75, 0xbc, 0xe2, 0x11, 0x48,
	0xbe, 0x86, 0xd6, 0x9c, 0x45, 0x0f, 0x98, 0xc5, 0xfe, 0xe1, 0xc0, 0x1a, 0x76, 0xc7, 0x47, 0x15,
	0xc5, 0x85, 0x81, 0x69, 0xe5, 0x3f, 0xa1, 0xe0, 0x3d, 0x96, 0x9a, 0xf4, 0xc1, 0x7e, 0xc0, 0xad,
	0x6e, 0x58, 0x87, 0xaa, 0x23, 0x79, 0x03, 0xcd, 0x0d, 0x4b, 0xd7, 0xe8, 0x37, 0x34, 0x59, 0xad,
	0xce, 0xe5, 0x36, 0x63, 0xab, 0x24, 0xba, 0x53, 0x3e, 0x6a, 0x42, 0x7e, 0x6c, 0xfc, 0x60, 0x9d,
	0xdc, 0xc0, 0x27, 0x4f, 0x54, 0xf9, 0x04, 0x71, 0xf0, 0x98, 0xb8, 0x57, 0x11, 0xab, 0xac, 0x3d,
	0xc2, 0x20, 0x81, 0x56, 0x59, 0x38, 0x21, 0xe0, 0xc8, 0x6d, 0x81, 0x25, 0x8b, 0x3e, 0x93, 0x6f,
	0xc0, 0x8d, 0xf2, 0xec, 0x3e, 0x59, 0x3c, 0x5b, 0x60, 0x19, 0x43, 0x5e, 0x42, 0xe7, 0xef, 0x9c,
	0x3f, 0x88, 0x82, 0x45, 0xe8, 0xdb, 0x9a, 0x66, 0x07, 0x04, 0x7f, 0x82, 0x6b, 0x1a, 0x4c, 0xbe,
	0x02, 0x97, 0x45, 0xb2, 0x9a, 0x5d, 0x6f, 0xec, 0x55, 0xac, 0xe7, 0x1a, 0xa5, 0xa5, 0x57, 0x7d,
	0x5d, 0x57, 0x5a, 0xcd, 0xf1, 0x47, 0xbe, 0x6e, 0x62, 0x82, 0x7f, 0x6d, 0x38, 0x7e, 0x7a, 0x3c,
	0xc9, 0xe7, 0xd0, 0x5d, 0xe5, 0xf1, 0x3a, 0xc5, 0xb0, 0x60, 0x72, 0x59, 0xde, 0x10, 0x0c, 0x34,
	0x63, 0x72, 0x49, 0x7e, 0x06, 0x67, 0x95, 0xc7, 0x46, 0x2d, 0x6f, 0xfc, 0xf6, 0xf9, 0x69, 0xaf,
	0xe1, 0x5f, 0xf3, 0x18, 0xa9, 0x4e, 0xac, 0xc5, 0xb3, 0xf7, 0xc4, 0x23, 0xe0, 0x64, 0x6c, 0x85,
	0xbe, 0x63, 0x30, 0x75, 0x26, 0x04, 0x6c, 0x21, 0xb9, 0xdf, 0x54, 0xd0, 0xf4, 0x80, 0x2a, 0x43,
	0x61, 0x49, 0x26, 0x7d, 0x77, 0x60, 0x0d, 0x6d, 0x85, 0x25, 0x99, 0x54, 0x15, 0xc7, 0x58, 0xe4,
	0x02, 0xe3, 0x50, 0x75, 0xb6, 0x65, 0x2a, 0x2e, 0xa1, 0x0f, 0xb8, 0x25, 0x27, 0xd0, 0xae, 0x46,
	0xd3, 0x6f, 0x6b, 0x6f, 0x6d, 0x2b, 0x7d, 0xcd, 0xd6, 0xf9, 0x1d, 0xdd, 0xb5, 0x5a, 0xdf, 0x72,
	0xdd, 0x4a, 0xaf, 0x7a, 0x44, 0x0a, 0x9e, 0x6c, 0x98, 0x44, 0x1f, 0x06, 0xd6, 0xb0, 0x47, 0x2b,
	0x93, 0x9c, 0xa9, 0x97, 0xe0, 0xaf, 0x75, 0xc2, 0x31, 0x0e, 0x39, 0x16, 0xa9, 0x6a, 0x68, 0x57,
	0xf7, 0xa0, 0x9e, 0x24, 0xa5, 0x9b, 0xda, 0x7b, 0x13, 0x45, 0x4d, 0x50, 0xf0, 0x25, 0xf4, 0xf6,
	0xd5, 0x21, 0x5d, 0x68, 0xad, 0x58, 0xc6, 0x16, 0x18, 0xf7, 0x0f, 0x48, 0x1b, 0x9c, 0x98, 0x49,
	0xd6, 0xb7, 0x2e, 0x3c, 0xe8, 0x25, 0xa5, 0xa6, 0xea, 0x7e, 0xc1, 0x12, 0x7a, 0xfb, 0x0f, 0x42,
	0x2d, 0x9d, 0xb5, 0x27, 0xdd, 0xee, 0x56, 0x8d, 0x67, 0x6f, 0xf5, 0x12, 0x3a, 0x02, 0x33, 0x91,
	0xc8, 0x64, 0x63, 0xfa, 0xd1, 0xa6, 0x3b, 0x20, 0x18, 0x42, 0x6f, 0x7f, 0x7a, 0x94, 0x06, 0x2b,
	0xb1, 0x28, 0x58, 0xf4, 0xa0, 0x3f, 0xd6, 0xa3, 0x95, 0x19, 0xbc, 0x02, 0x47, 0x6d, 0x0b, 0x39,
	0x06, 0x57, 0x2c, 0xd9, 0xf8, 0xbb, 0xef, 0xcb, 0x80, 0xd2, 0x0a, 0xfe, 0xb1, 0xc0, 0xd1, 0xc3,
	0xf3, 0x1a, 0x9a, 0x42, 0x62, 0x21, 0x7c, 0x4b, 0x2b, 0xf4, 0x62, 0x5f, 0xa1, 0xd1, 0x6f, 0x12,
	0x0b, 0x6a, 0xfc, 0x27, 0x12, 0x1c, 0x65, 0x92, 0xd7, 0xe0, 0x31, 0x29, 0x79, 0x32, 0x5f, 0x4b,
	0x0c, 0x77, 0xf7, 0x9c, 0x1e, 0xd0, 0xc3, 0x1a, 0xbf, 0x56, 0x57, 0x3e, 0x83, 0x2e, 0xa6, 0xb8,
	0xc2, 0x4c, 0xea, 0x29, 0x78, 0x66, 0x07, 0xa7, 0x07, 0x14, 0xca, 0xd0, 0x0f, 0xb8, 0xbd, 0x00,
	0x68, 0x0b, 0x4c, 0x31, 0x92, 0x39, 0x7f, 0x53, 0x80, 0x6b, 0xf6, 0x4a, 0xe9, 0x7f, 0x7d, 0x73,
	0x33, 0xeb, 0x1f, 0x10, 0x00, 0xf7, 0x3d, 0x9d, 0x9c, 0xdf, 0x4e, 0xfa, 0x96, 0x42, 0xe9, 0xe4,
	0xfc, 0xb2, 0xdf, 0x50, 0xe8, 0xef, 0xb3, 0x4b, 0x85, 0xda, 0xea, 0x7c, 0x39, 0xf9, 0x65, 0x72,
	0x3b, 0xe9, 0x37, 0xc9, 0x31, 0x10, 0x73, 0x0e, 0x6f, 0xa7, 0x93, 0xeb, 0xb0, 0xcc, 0x74, 0x15,
	0x6e, 0xce, 0x06, 0x2f, 0xe3, 0x5b, 0x17, 0xef, 0xfe, 0x38, 0x5b, 0x24, 0x72, 0xb9, 0x9e, 0x8f,
	0xa2, 0x7c, 0x75, 0xaa, 0x5e, 0xdc, 0x24, 0xca, 0x79, 0x71, 0x5a, 0x3f, 0xcc, 0xa7, 0xaa, 0x7e,
	0x71, 0x9a, 0x64, 0x12, 0x79, 0xc6, 0x52, 0x6d, 0xea, 0x3f, 0xba, 0xb9, 0xab, 0x7f, 0xbe, 0xfd,
	0x3f, 0x00, 0x00, 0xff, 0xff, 0x30, 0x3e, 0x4e, 0x33, 0x01, 0x07, 0x00, 0x00,
}
