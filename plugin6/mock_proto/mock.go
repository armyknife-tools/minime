// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hashicorp/terraform/internal/tfplugin6 (interfaces: ProviderClient)

// Package mock_tfplugin6 is a generated GoMock package.
package mock_tfplugin6

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	tfplugin6 "github.com/hashicorp/terraform/internal/tfplugin6"
	grpc "google.golang.org/grpc"
	reflect "reflect"
)

// MockProviderClient is a mock of ProviderClient interface
type MockProviderClient struct {
	ctrl     *gomock.Controller
	recorder *MockProviderClientMockRecorder
}

// MockProviderClientMockRecorder is the mock recorder for MockProviderClient
type MockProviderClientMockRecorder struct {
	mock *MockProviderClient
}

// NewMockProviderClient creates a new mock instance
func NewMockProviderClient(ctrl *gomock.Controller) *MockProviderClient {
	mock := &MockProviderClient{ctrl: ctrl}
	mock.recorder = &MockProviderClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProviderClient) EXPECT() *MockProviderClientMockRecorder {
	return m.recorder
}

// ApplyResourceChange mocks base method
func (m *MockProviderClient) ApplyResourceChange(arg0 context.Context, arg1 *tfplugin6.ApplyResourceChange_Request, arg2 ...grpc.CallOption) (*tfplugin6.ApplyResourceChange_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ApplyResourceChange", varargs...)
	ret0, _ := ret[0].(*tfplugin6.ApplyResourceChange_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApplyResourceChange indicates an expected call of ApplyResourceChange
func (mr *MockProviderClientMockRecorder) ApplyResourceChange(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplyResourceChange", reflect.TypeOf((*MockProviderClient)(nil).ApplyResourceChange), varargs...)
}

// ConfigureProvider mocks base method
func (m *MockProviderClient) ConfigureProvider(arg0 context.Context, arg1 *tfplugin6.ConfigureProvider_Request, arg2 ...grpc.CallOption) (*tfplugin6.ConfigureProvider_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ConfigureProvider", varargs...)
	ret0, _ := ret[0].(*tfplugin6.ConfigureProvider_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConfigureProvider indicates an expected call of ConfigureProvider
func (mr *MockProviderClientMockRecorder) ConfigureProvider(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfigureProvider", reflect.TypeOf((*MockProviderClient)(nil).ConfigureProvider), varargs...)
}

// GetProviderSchema mocks base method
func (m *MockProviderClient) GetProviderSchema(arg0 context.Context, arg1 *tfplugin6.GetProviderSchema_Request, arg2 ...grpc.CallOption) (*tfplugin6.GetProviderSchema_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetProviderSchema", varargs...)
	ret0, _ := ret[0].(*tfplugin6.GetProviderSchema_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProviderSchema indicates an expected call of GetProviderSchema
func (mr *MockProviderClientMockRecorder) GetProviderSchema(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProviderSchema", reflect.TypeOf((*MockProviderClient)(nil).GetProviderSchema), varargs...)
}

// ImportResourceState mocks base method
func (m *MockProviderClient) ImportResourceState(arg0 context.Context, arg1 *tfplugin6.ImportResourceState_Request, arg2 ...grpc.CallOption) (*tfplugin6.ImportResourceState_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ImportResourceState", varargs...)
	ret0, _ := ret[0].(*tfplugin6.ImportResourceState_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImportResourceState indicates an expected call of ImportResourceState
func (mr *MockProviderClientMockRecorder) ImportResourceState(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImportResourceState", reflect.TypeOf((*MockProviderClient)(nil).ImportResourceState), varargs...)
}

// PlanResourceChange mocks base method
func (m *MockProviderClient) PlanResourceChange(arg0 context.Context, arg1 *tfplugin6.PlanResourceChange_Request, arg2 ...grpc.CallOption) (*tfplugin6.PlanResourceChange_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PlanResourceChange", varargs...)
	ret0, _ := ret[0].(*tfplugin6.PlanResourceChange_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PlanResourceChange indicates an expected call of PlanResourceChange
func (mr *MockProviderClientMockRecorder) PlanResourceChange(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlanResourceChange", reflect.TypeOf((*MockProviderClient)(nil).PlanResourceChange), varargs...)
}

// ReadDataSource mocks base method
func (m *MockProviderClient) ReadDataSource(arg0 context.Context, arg1 *tfplugin6.ReadDataSource_Request, arg2 ...grpc.CallOption) (*tfplugin6.ReadDataSource_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReadDataSource", varargs...)
	ret0, _ := ret[0].(*tfplugin6.ReadDataSource_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadDataSource indicates an expected call of ReadDataSource
func (mr *MockProviderClientMockRecorder) ReadDataSource(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadDataSource", reflect.TypeOf((*MockProviderClient)(nil).ReadDataSource), varargs...)
}

// ReadResource mocks base method
func (m *MockProviderClient) ReadResource(arg0 context.Context, arg1 *tfplugin6.ReadResource_Request, arg2 ...grpc.CallOption) (*tfplugin6.ReadResource_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReadResource", varargs...)
	ret0, _ := ret[0].(*tfplugin6.ReadResource_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadResource indicates an expected call of ReadResource
func (mr *MockProviderClientMockRecorder) ReadResource(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadResource", reflect.TypeOf((*MockProviderClient)(nil).ReadResource), varargs...)
}

// StopProvider mocks base method
func (m *MockProviderClient) StopProvider(arg0 context.Context, arg1 *tfplugin6.StopProvider_Request, arg2 ...grpc.CallOption) (*tfplugin6.StopProvider_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "StopProvider", varargs...)
	ret0, _ := ret[0].(*tfplugin6.StopProvider_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StopProvider indicates an expected call of StopProvider
func (mr *MockProviderClientMockRecorder) StopProvider(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopProvider", reflect.TypeOf((*MockProviderClient)(nil).StopProvider), varargs...)
}

// UpgradeResourceState mocks base method
func (m *MockProviderClient) UpgradeResourceState(arg0 context.Context, arg1 *tfplugin6.UpgradeResourceState_Request, arg2 ...grpc.CallOption) (*tfplugin6.UpgradeResourceState_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpgradeResourceState", varargs...)
	ret0, _ := ret[0].(*tfplugin6.UpgradeResourceState_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpgradeResourceState indicates an expected call of UpgradeResourceState
func (mr *MockProviderClientMockRecorder) UpgradeResourceState(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpgradeResourceState", reflect.TypeOf((*MockProviderClient)(nil).UpgradeResourceState), varargs...)
}

// ValidateDataSourceConfig mocks base method
func (m *MockProviderClient) ValidateDataSourceConfig(arg0 context.Context, arg1 *tfplugin6.ValidateDataSourceConfig_Request, arg2 ...grpc.CallOption) (*tfplugin6.ValidateDataSourceConfig_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateDataSourceConfig", varargs...)
	ret0, _ := ret[0].(*tfplugin6.ValidateDataSourceConfig_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateDataSourceConfig indicates an expected call of ValidateDataSourceConfig
func (mr *MockProviderClientMockRecorder) ValidateDataSourceConfig(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateDataSourceConfig", reflect.TypeOf((*MockProviderClient)(nil).ValidateDataSourceConfig), varargs...)
}

// ValidateProviderConfig mocks base method
func (m *MockProviderClient) ValidateProviderConfig(arg0 context.Context, arg1 *tfplugin6.ValidateProviderConfig_Request, arg2 ...grpc.CallOption) (*tfplugin6.ValidateProviderConfig_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateProviderConfig", varargs...)
	ret0, _ := ret[0].(*tfplugin6.ValidateProviderConfig_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateProviderConfig indicates an expected call of ValidateProviderConfig
func (mr *MockProviderClientMockRecorder) ValidateProviderConfig(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateProviderConfig", reflect.TypeOf((*MockProviderClient)(nil).ValidateProviderConfig), varargs...)
}

// ValidateResourceConfig mocks base method
func (m *MockProviderClient) ValidateResourceConfig(arg0 context.Context, arg1 *tfplugin6.ValidateResourceConfig_Request, arg2 ...grpc.CallOption) (*tfplugin6.ValidateResourceConfig_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateResourceConfig", varargs...)
	ret0, _ := ret[0].(*tfplugin6.ValidateResourceConfig_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateResourceConfig indicates an expected call of ValidateResourceConfig
func (mr *MockProviderClientMockRecorder) ValidateResourceConfig(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateResourceConfig", reflect.TypeOf((*MockProviderClient)(nil).ValidateResourceConfig), varargs...)
}
