// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/opentofu/opentofu/internal/tfplugin6 (interfaces: ProviderClient)
//
// Generated by this command:
//
//	mockgen -destination mock.go github.com/opentofu/opentofu/internal/tfplugin6 ProviderClient
//

// Package mock_tfplugin6 is a generated GoMock package.
package mock_tfplugin6

import (
	context "context"
	reflect "reflect"

	tfplugin6 "github.com/opentofu/opentofu/internal/tfplugin6"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockProviderClient is a mock of ProviderClient interface.
type MockProviderClient struct {
	ctrl     *gomock.Controller
	recorder *MockProviderClientMockRecorder
}

// MockProviderClientMockRecorder is the mock recorder for MockProviderClient.
type MockProviderClientMockRecorder struct {
	mock *MockProviderClient
}

// NewMockProviderClient creates a new mock instance.
func NewMockProviderClient(ctrl *gomock.Controller) *MockProviderClient {
	mock := &MockProviderClient{ctrl: ctrl}
	mock.recorder = &MockProviderClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProviderClient) EXPECT() *MockProviderClientMockRecorder {
	return m.recorder
}

// ApplyResourceChange mocks base method.
func (m *MockProviderClient) ApplyResourceChange(arg0 context.Context, arg1 *tfplugin6.ApplyResourceChange_Request, arg2 ...grpc.CallOption) (*tfplugin6.ApplyResourceChange_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ApplyResourceChange", varargs...)
	ret0, _ := ret[0].(*tfplugin6.ApplyResourceChange_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApplyResourceChange indicates an expected call of ApplyResourceChange.
func (mr *MockProviderClientMockRecorder) ApplyResourceChange(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplyResourceChange", reflect.TypeOf((*MockProviderClient)(nil).ApplyResourceChange), varargs...)
}

// CallFunction mocks base method.
func (m *MockProviderClient) CallFunction(arg0 context.Context, arg1 *tfplugin6.CallFunction_Request, arg2 ...grpc.CallOption) (*tfplugin6.CallFunction_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CallFunction", varargs...)
	ret0, _ := ret[0].(*tfplugin6.CallFunction_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CallFunction indicates an expected call of CallFunction.
func (mr *MockProviderClientMockRecorder) CallFunction(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CallFunction", reflect.TypeOf((*MockProviderClient)(nil).CallFunction), varargs...)
}

// CloseEphemeralResource mocks base method.
func (m *MockProviderClient) CloseEphemeralResource(arg0 context.Context, arg1 *tfplugin6.CloseEphemeralResource_Request, arg2 ...grpc.CallOption) (*tfplugin6.CloseEphemeralResource_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CloseEphemeralResource", varargs...)
	ret0, _ := ret[0].(*tfplugin6.CloseEphemeralResource_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloseEphemeralResource indicates an expected call of CloseEphemeralResource.
func (mr *MockProviderClientMockRecorder) CloseEphemeralResource(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseEphemeralResource", reflect.TypeOf((*MockProviderClient)(nil).CloseEphemeralResource), varargs...)
}

// ConfigureProvider mocks base method.
func (m *MockProviderClient) ConfigureProvider(arg0 context.Context, arg1 *tfplugin6.ConfigureProvider_Request, arg2 ...grpc.CallOption) (*tfplugin6.ConfigureProvider_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ConfigureProvider", varargs...)
	ret0, _ := ret[0].(*tfplugin6.ConfigureProvider_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConfigureProvider indicates an expected call of ConfigureProvider.
func (mr *MockProviderClientMockRecorder) ConfigureProvider(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfigureProvider", reflect.TypeOf((*MockProviderClient)(nil).ConfigureProvider), varargs...)
}

// GetFunctions mocks base method.
func (m *MockProviderClient) GetFunctions(arg0 context.Context, arg1 *tfplugin6.GetFunctions_Request, arg2 ...grpc.CallOption) (*tfplugin6.GetFunctions_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetFunctions", varargs...)
	ret0, _ := ret[0].(*tfplugin6.GetFunctions_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFunctions indicates an expected call of GetFunctions.
func (mr *MockProviderClientMockRecorder) GetFunctions(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFunctions", reflect.TypeOf((*MockProviderClient)(nil).GetFunctions), varargs...)
}

// GetMetadata mocks base method.
func (m *MockProviderClient) GetMetadata(arg0 context.Context, arg1 *tfplugin6.GetMetadata_Request, arg2 ...grpc.CallOption) (*tfplugin6.GetMetadata_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetMetadata", varargs...)
	ret0, _ := ret[0].(*tfplugin6.GetMetadata_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetadata indicates an expected call of GetMetadata.
func (mr *MockProviderClientMockRecorder) GetMetadata(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetadata", reflect.TypeOf((*MockProviderClient)(nil).GetMetadata), varargs...)
}

// GetProviderSchema mocks base method.
func (m *MockProviderClient) GetProviderSchema(arg0 context.Context, arg1 *tfplugin6.GetProviderSchema_Request, arg2 ...grpc.CallOption) (*tfplugin6.GetProviderSchema_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetProviderSchema", varargs...)
	ret0, _ := ret[0].(*tfplugin6.GetProviderSchema_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProviderSchema indicates an expected call of GetProviderSchema.
func (mr *MockProviderClientMockRecorder) GetProviderSchema(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProviderSchema", reflect.TypeOf((*MockProviderClient)(nil).GetProviderSchema), varargs...)
}

// ImportResourceState mocks base method.
func (m *MockProviderClient) ImportResourceState(arg0 context.Context, arg1 *tfplugin6.ImportResourceState_Request, arg2 ...grpc.CallOption) (*tfplugin6.ImportResourceState_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ImportResourceState", varargs...)
	ret0, _ := ret[0].(*tfplugin6.ImportResourceState_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImportResourceState indicates an expected call of ImportResourceState.
func (mr *MockProviderClientMockRecorder) ImportResourceState(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImportResourceState", reflect.TypeOf((*MockProviderClient)(nil).ImportResourceState), varargs...)
}

// MoveResourceState mocks base method.
func (m *MockProviderClient) MoveResourceState(arg0 context.Context, arg1 *tfplugin6.MoveResourceState_Request, arg2 ...grpc.CallOption) (*tfplugin6.MoveResourceState_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "MoveResourceState", varargs...)
	ret0, _ := ret[0].(*tfplugin6.MoveResourceState_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MoveResourceState indicates an expected call of MoveResourceState.
func (mr *MockProviderClientMockRecorder) MoveResourceState(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MoveResourceState", reflect.TypeOf((*MockProviderClient)(nil).MoveResourceState), varargs...)
}

// OpenEphemeralResource mocks base method.
func (m *MockProviderClient) OpenEphemeralResource(arg0 context.Context, arg1 *tfplugin6.OpenEphemeralResource_Request, arg2 ...grpc.CallOption) (*tfplugin6.OpenEphemeralResource_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "OpenEphemeralResource", varargs...)
	ret0, _ := ret[0].(*tfplugin6.OpenEphemeralResource_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenEphemeralResource indicates an expected call of OpenEphemeralResource.
func (mr *MockProviderClientMockRecorder) OpenEphemeralResource(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenEphemeralResource", reflect.TypeOf((*MockProviderClient)(nil).OpenEphemeralResource), varargs...)
}

// PlanResourceChange mocks base method.
func (m *MockProviderClient) PlanResourceChange(arg0 context.Context, arg1 *tfplugin6.PlanResourceChange_Request, arg2 ...grpc.CallOption) (*tfplugin6.PlanResourceChange_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PlanResourceChange", varargs...)
	ret0, _ := ret[0].(*tfplugin6.PlanResourceChange_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PlanResourceChange indicates an expected call of PlanResourceChange.
func (mr *MockProviderClientMockRecorder) PlanResourceChange(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlanResourceChange", reflect.TypeOf((*MockProviderClient)(nil).PlanResourceChange), varargs...)
}

// ReadDataSource mocks base method.
func (m *MockProviderClient) ReadDataSource(arg0 context.Context, arg1 *tfplugin6.ReadDataSource_Request, arg2 ...grpc.CallOption) (*tfplugin6.ReadDataSource_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReadDataSource", varargs...)
	ret0, _ := ret[0].(*tfplugin6.ReadDataSource_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadDataSource indicates an expected call of ReadDataSource.
func (mr *MockProviderClientMockRecorder) ReadDataSource(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadDataSource", reflect.TypeOf((*MockProviderClient)(nil).ReadDataSource), varargs...)
}

// ReadResource mocks base method.
func (m *MockProviderClient) ReadResource(arg0 context.Context, arg1 *tfplugin6.ReadResource_Request, arg2 ...grpc.CallOption) (*tfplugin6.ReadResource_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReadResource", varargs...)
	ret0, _ := ret[0].(*tfplugin6.ReadResource_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadResource indicates an expected call of ReadResource.
func (mr *MockProviderClientMockRecorder) ReadResource(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadResource", reflect.TypeOf((*MockProviderClient)(nil).ReadResource), varargs...)
}

// RenewEphemeralResource mocks base method.
func (m *MockProviderClient) RenewEphemeralResource(arg0 context.Context, arg1 *tfplugin6.RenewEphemeralResource_Request, arg2 ...grpc.CallOption) (*tfplugin6.RenewEphemeralResource_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RenewEphemeralResource", varargs...)
	ret0, _ := ret[0].(*tfplugin6.RenewEphemeralResource_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RenewEphemeralResource indicates an expected call of RenewEphemeralResource.
func (mr *MockProviderClientMockRecorder) RenewEphemeralResource(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenewEphemeralResource", reflect.TypeOf((*MockProviderClient)(nil).RenewEphemeralResource), varargs...)
}

// StopProvider mocks base method.
func (m *MockProviderClient) StopProvider(arg0 context.Context, arg1 *tfplugin6.StopProvider_Request, arg2 ...grpc.CallOption) (*tfplugin6.StopProvider_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "StopProvider", varargs...)
	ret0, _ := ret[0].(*tfplugin6.StopProvider_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StopProvider indicates an expected call of StopProvider.
func (mr *MockProviderClientMockRecorder) StopProvider(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopProvider", reflect.TypeOf((*MockProviderClient)(nil).StopProvider), varargs...)
}

// UpgradeResourceState mocks base method.
func (m *MockProviderClient) UpgradeResourceState(arg0 context.Context, arg1 *tfplugin6.UpgradeResourceState_Request, arg2 ...grpc.CallOption) (*tfplugin6.UpgradeResourceState_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpgradeResourceState", varargs...)
	ret0, _ := ret[0].(*tfplugin6.UpgradeResourceState_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpgradeResourceState indicates an expected call of UpgradeResourceState.
func (mr *MockProviderClientMockRecorder) UpgradeResourceState(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpgradeResourceState", reflect.TypeOf((*MockProviderClient)(nil).UpgradeResourceState), varargs...)
}

// ValidateDataResourceConfig mocks base method.
func (m *MockProviderClient) ValidateDataResourceConfig(arg0 context.Context, arg1 *tfplugin6.ValidateDataResourceConfig_Request, arg2 ...grpc.CallOption) (*tfplugin6.ValidateDataResourceConfig_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateDataResourceConfig", varargs...)
	ret0, _ := ret[0].(*tfplugin6.ValidateDataResourceConfig_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateDataResourceConfig indicates an expected call of ValidateDataResourceConfig.
func (mr *MockProviderClientMockRecorder) ValidateDataResourceConfig(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateDataResourceConfig", reflect.TypeOf((*MockProviderClient)(nil).ValidateDataResourceConfig), varargs...)
}

// ValidateEphemeralResourceConfig mocks base method.
func (m *MockProviderClient) ValidateEphemeralResourceConfig(arg0 context.Context, arg1 *tfplugin6.ValidateEphemeralResourceConfig_Request, arg2 ...grpc.CallOption) (*tfplugin6.ValidateEphemeralResourceConfig_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateEphemeralResourceConfig", varargs...)
	ret0, _ := ret[0].(*tfplugin6.ValidateEphemeralResourceConfig_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateEphemeralResourceConfig indicates an expected call of ValidateEphemeralResourceConfig.
func (mr *MockProviderClientMockRecorder) ValidateEphemeralResourceConfig(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateEphemeralResourceConfig", reflect.TypeOf((*MockProviderClient)(nil).ValidateEphemeralResourceConfig), varargs...)
}

// ValidateProviderConfig mocks base method.
func (m *MockProviderClient) ValidateProviderConfig(arg0 context.Context, arg1 *tfplugin6.ValidateProviderConfig_Request, arg2 ...grpc.CallOption) (*tfplugin6.ValidateProviderConfig_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateProviderConfig", varargs...)
	ret0, _ := ret[0].(*tfplugin6.ValidateProviderConfig_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateProviderConfig indicates an expected call of ValidateProviderConfig.
func (mr *MockProviderClientMockRecorder) ValidateProviderConfig(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateProviderConfig", reflect.TypeOf((*MockProviderClient)(nil).ValidateProviderConfig), varargs...)
}

// ValidateResourceConfig mocks base method.
func (m *MockProviderClient) ValidateResourceConfig(arg0 context.Context, arg1 *tfplugin6.ValidateResourceConfig_Request, arg2 ...grpc.CallOption) (*tfplugin6.ValidateResourceConfig_Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateResourceConfig", varargs...)
	ret0, _ := ret[0].(*tfplugin6.ValidateResourceConfig_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateResourceConfig indicates an expected call of ValidateResourceConfig.
func (mr *MockProviderClientMockRecorder) ValidateResourceConfig(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateResourceConfig", reflect.TypeOf((*MockProviderClient)(nil).ValidateResourceConfig), varargs...)
}
