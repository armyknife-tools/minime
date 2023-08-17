// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/placeholderplaceholderplaceholder/opentf/internal/tfplugin5 (interfaces: ProviderClient,ProvisionerClient,Provisioner_ProvisionResourceClient,Provisioner_ProvisionResourceServer)

// Package mock_tfplugin5 is a generated GoMock package.
package mock_tfplugin5

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	tfplugin5 "github.com/placeholderplaceholderplaceholder/opentf/internal/tfplugin5"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
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
func (m *MockProviderClient) ApplyResourceChange(arg0 context.Context, arg1 *tfplugin5.ApplyResourceChange_Request, arg2 ...grpc.CallOption) (*tfplugin5.ApplyResourceChange_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ApplyResourceChange", varargs...)
	ret0, _ := ret[0].(*tfplugin5.ApplyResourceChange_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApplyResourceChange indicates an expected call of ApplyResourceChange.
func (mr *MockProviderClientMockRecorder) ApplyResourceChange(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplyResourceChange", reflect.TypeOf((*MockProviderClient)(nil).ApplyResourceChange), varargs...)
}

// Configure mocks base method.
func (m *MockProviderClient) Configure(arg0 context.Context, arg1 *tfplugin5.Configure_Request, arg2 ...grpc.CallOption) (*tfplugin5.Configure_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Configure", varargs...)
	ret0, _ := ret[0].(*tfplugin5.Configure_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Configure indicates an expected call of Configure.
func (mr *MockProviderClientMockRecorder) Configure(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Configure", reflect.TypeOf((*MockProviderClient)(nil).Configure), varargs...)
}

// GetSchema mocks base method.
func (m *MockProviderClient) GetSchema(arg0 context.Context, arg1 *tfplugin5.GetProviderSchema_Request, arg2 ...grpc.CallOption) (*tfplugin5.GetProviderSchema_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSchema", varargs...)
	ret0, _ := ret[0].(*tfplugin5.GetProviderSchema_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSchema indicates an expected call of GetSchema.
func (mr *MockProviderClientMockRecorder) GetSchema(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSchema", reflect.TypeOf((*MockProviderClient)(nil).GetSchema), varargs...)
}

// ImportResourceState mocks base method.
func (m *MockProviderClient) ImportResourceState(arg0 context.Context, arg1 *tfplugin5.ImportResourceState_Request, arg2 ...grpc.CallOption) (*tfplugin5.ImportResourceState_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ImportResourceState", varargs...)
	ret0, _ := ret[0].(*tfplugin5.ImportResourceState_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImportResourceState indicates an expected call of ImportResourceState.
func (mr *MockProviderClientMockRecorder) ImportResourceState(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImportResourceState", reflect.TypeOf((*MockProviderClient)(nil).ImportResourceState), varargs...)
}

// PlanResourceChange mocks base method.
func (m *MockProviderClient) PlanResourceChange(arg0 context.Context, arg1 *tfplugin5.PlanResourceChange_Request, arg2 ...grpc.CallOption) (*tfplugin5.PlanResourceChange_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PlanResourceChange", varargs...)
	ret0, _ := ret[0].(*tfplugin5.PlanResourceChange_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PlanResourceChange indicates an expected call of PlanResourceChange.
func (mr *MockProviderClientMockRecorder) PlanResourceChange(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlanResourceChange", reflect.TypeOf((*MockProviderClient)(nil).PlanResourceChange), varargs...)
}

// PrepareProviderConfig mocks base method.
func (m *MockProviderClient) PrepareProviderConfig(arg0 context.Context, arg1 *tfplugin5.PrepareProviderConfig_Request, arg2 ...grpc.CallOption) (*tfplugin5.PrepareProviderConfig_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PrepareProviderConfig", varargs...)
	ret0, _ := ret[0].(*tfplugin5.PrepareProviderConfig_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrepareProviderConfig indicates an expected call of PrepareProviderConfig.
func (mr *MockProviderClientMockRecorder) PrepareProviderConfig(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareProviderConfig", reflect.TypeOf((*MockProviderClient)(nil).PrepareProviderConfig), varargs...)
}

// ReadDataSource mocks base method.
func (m *MockProviderClient) ReadDataSource(arg0 context.Context, arg1 *tfplugin5.ReadDataSource_Request, arg2 ...grpc.CallOption) (*tfplugin5.ReadDataSource_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReadDataSource", varargs...)
	ret0, _ := ret[0].(*tfplugin5.ReadDataSource_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadDataSource indicates an expected call of ReadDataSource.
func (mr *MockProviderClientMockRecorder) ReadDataSource(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadDataSource", reflect.TypeOf((*MockProviderClient)(nil).ReadDataSource), varargs...)
}

// ReadResource mocks base method.
func (m *MockProviderClient) ReadResource(arg0 context.Context, arg1 *tfplugin5.ReadResource_Request, arg2 ...grpc.CallOption) (*tfplugin5.ReadResource_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReadResource", varargs...)
	ret0, _ := ret[0].(*tfplugin5.ReadResource_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadResource indicates an expected call of ReadResource.
func (mr *MockProviderClientMockRecorder) ReadResource(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadResource", reflect.TypeOf((*MockProviderClient)(nil).ReadResource), varargs...)
}

// Stop mocks base method.
func (m *MockProviderClient) Stop(arg0 context.Context, arg1 *tfplugin5.Stop_Request, arg2 ...grpc.CallOption) (*tfplugin5.Stop_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Stop", varargs...)
	ret0, _ := ret[0].(*tfplugin5.Stop_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stop indicates an expected call of Stop.
func (mr *MockProviderClientMockRecorder) Stop(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockProviderClient)(nil).Stop), varargs...)
}

// UpgradeResourceState mocks base method.
func (m *MockProviderClient) UpgradeResourceState(arg0 context.Context, arg1 *tfplugin5.UpgradeResourceState_Request, arg2 ...grpc.CallOption) (*tfplugin5.UpgradeResourceState_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpgradeResourceState", varargs...)
	ret0, _ := ret[0].(*tfplugin5.UpgradeResourceState_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpgradeResourceState indicates an expected call of UpgradeResourceState.
func (mr *MockProviderClientMockRecorder) UpgradeResourceState(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpgradeResourceState", reflect.TypeOf((*MockProviderClient)(nil).UpgradeResourceState), varargs...)
}

// ValidateDataSourceConfig mocks base method.
func (m *MockProviderClient) ValidateDataSourceConfig(arg0 context.Context, arg1 *tfplugin5.ValidateDataSourceConfig_Request, arg2 ...grpc.CallOption) (*tfplugin5.ValidateDataSourceConfig_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateDataSourceConfig", varargs...)
	ret0, _ := ret[0].(*tfplugin5.ValidateDataSourceConfig_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateDataSourceConfig indicates an expected call of ValidateDataSourceConfig.
func (mr *MockProviderClientMockRecorder) ValidateDataSourceConfig(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateDataSourceConfig", reflect.TypeOf((*MockProviderClient)(nil).ValidateDataSourceConfig), varargs...)
}

// ValidateResourceTypeConfig mocks base method.
func (m *MockProviderClient) ValidateResourceTypeConfig(arg0 context.Context, arg1 *tfplugin5.ValidateResourceTypeConfig_Request, arg2 ...grpc.CallOption) (*tfplugin5.ValidateResourceTypeConfig_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateResourceTypeConfig", varargs...)
	ret0, _ := ret[0].(*tfplugin5.ValidateResourceTypeConfig_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateResourceTypeConfig indicates an expected call of ValidateResourceTypeConfig.
func (mr *MockProviderClientMockRecorder) ValidateResourceTypeConfig(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateResourceTypeConfig", reflect.TypeOf((*MockProviderClient)(nil).ValidateResourceTypeConfig), varargs...)
}

// MockProvisionerClient is a mock of ProvisionerClient interface.
type MockProvisionerClient struct {
	ctrl     *gomock.Controller
	recorder *MockProvisionerClientMockRecorder
}

// MockProvisionerClientMockRecorder is the mock recorder for MockProvisionerClient.
type MockProvisionerClientMockRecorder struct {
	mock *MockProvisionerClient
}

// NewMockProvisionerClient creates a new mock instance.
func NewMockProvisionerClient(ctrl *gomock.Controller) *MockProvisionerClient {
	mock := &MockProvisionerClient{ctrl: ctrl}
	mock.recorder = &MockProvisionerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProvisionerClient) EXPECT() *MockProvisionerClientMockRecorder {
	return m.recorder
}

// GetSchema mocks base method.
func (m *MockProvisionerClient) GetSchema(arg0 context.Context, arg1 *tfplugin5.GetProvisionerSchema_Request, arg2 ...grpc.CallOption) (*tfplugin5.GetProvisionerSchema_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSchema", varargs...)
	ret0, _ := ret[0].(*tfplugin5.GetProvisionerSchema_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSchema indicates an expected call of GetSchema.
func (mr *MockProvisionerClientMockRecorder) GetSchema(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSchema", reflect.TypeOf((*MockProvisionerClient)(nil).GetSchema), varargs...)
}

// ProvisionResource mocks base method.
func (m *MockProvisionerClient) ProvisionResource(arg0 context.Context, arg1 *tfplugin5.ProvisionResource_Request, arg2 ...grpc.CallOption) (tfplugin5.Provisioner_ProvisionResourceClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ProvisionResource", varargs...)
	ret0, _ := ret[0].(tfplugin5.Provisioner_ProvisionResourceClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProvisionResource indicates an expected call of ProvisionResource.
func (mr *MockProvisionerClientMockRecorder) ProvisionResource(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProvisionResource", reflect.TypeOf((*MockProvisionerClient)(nil).ProvisionResource), varargs...)
}

// Stop mocks base method.
func (m *MockProvisionerClient) Stop(arg0 context.Context, arg1 *tfplugin5.Stop_Request, arg2 ...grpc.CallOption) (*tfplugin5.Stop_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Stop", varargs...)
	ret0, _ := ret[0].(*tfplugin5.Stop_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stop indicates an expected call of Stop.
func (mr *MockProvisionerClientMockRecorder) Stop(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockProvisionerClient)(nil).Stop), varargs...)
}

// ValidateProvisionerConfig mocks base method.
func (m *MockProvisionerClient) ValidateProvisionerConfig(arg0 context.Context, arg1 *tfplugin5.ValidateProvisionerConfig_Request, arg2 ...grpc.CallOption) (*tfplugin5.ValidateProvisionerConfig_Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateProvisionerConfig", varargs...)
	ret0, _ := ret[0].(*tfplugin5.ValidateProvisionerConfig_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateProvisionerConfig indicates an expected call of ValidateProvisionerConfig.
func (mr *MockProvisionerClientMockRecorder) ValidateProvisionerConfig(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateProvisionerConfig", reflect.TypeOf((*MockProvisionerClient)(nil).ValidateProvisionerConfig), varargs...)
}

// MockProvisioner_ProvisionResourceClient is a mock of Provisioner_ProvisionResourceClient interface.
type MockProvisioner_ProvisionResourceClient struct {
	ctrl     *gomock.Controller
	recorder *MockProvisioner_ProvisionResourceClientMockRecorder
}

// MockProvisioner_ProvisionResourceClientMockRecorder is the mock recorder for MockProvisioner_ProvisionResourceClient.
type MockProvisioner_ProvisionResourceClientMockRecorder struct {
	mock *MockProvisioner_ProvisionResourceClient
}

// NewMockProvisioner_ProvisionResourceClient creates a new mock instance.
func NewMockProvisioner_ProvisionResourceClient(ctrl *gomock.Controller) *MockProvisioner_ProvisionResourceClient {
	mock := &MockProvisioner_ProvisionResourceClient{ctrl: ctrl}
	mock.recorder = &MockProvisioner_ProvisionResourceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProvisioner_ProvisionResourceClient) EXPECT() *MockProvisioner_ProvisionResourceClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockProvisioner_ProvisionResourceClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockProvisioner_ProvisionResourceClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockProvisioner_ProvisionResourceClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockProvisioner_ProvisionResourceClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockProvisioner_ProvisionResourceClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockProvisioner_ProvisionResourceClient)(nil).Context))
}

// Header mocks base method.
func (m *MockProvisioner_ProvisionResourceClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockProvisioner_ProvisionResourceClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockProvisioner_ProvisionResourceClient)(nil).Header))
}

// Recv mocks base method.
func (m *MockProvisioner_ProvisionResourceClient) Recv() (*tfplugin5.ProvisionResource_Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*tfplugin5.ProvisionResource_Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockProvisioner_ProvisionResourceClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockProvisioner_ProvisionResourceClient)(nil).Recv))
}

// RecvMsg mocks base method.
func (m *MockProvisioner_ProvisionResourceClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockProvisioner_ProvisionResourceClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockProvisioner_ProvisionResourceClient)(nil).RecvMsg), arg0)
}

// SendMsg mocks base method.
func (m *MockProvisioner_ProvisionResourceClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockProvisioner_ProvisionResourceClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockProvisioner_ProvisionResourceClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method.
func (m *MockProvisioner_ProvisionResourceClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockProvisioner_ProvisionResourceClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockProvisioner_ProvisionResourceClient)(nil).Trailer))
}

// MockProvisioner_ProvisionResourceServer is a mock of Provisioner_ProvisionResourceServer interface.
type MockProvisioner_ProvisionResourceServer struct {
	ctrl     *gomock.Controller
	recorder *MockProvisioner_ProvisionResourceServerMockRecorder
}

// MockProvisioner_ProvisionResourceServerMockRecorder is the mock recorder for MockProvisioner_ProvisionResourceServer.
type MockProvisioner_ProvisionResourceServerMockRecorder struct {
	mock *MockProvisioner_ProvisionResourceServer
}

// NewMockProvisioner_ProvisionResourceServer creates a new mock instance.
func NewMockProvisioner_ProvisionResourceServer(ctrl *gomock.Controller) *MockProvisioner_ProvisionResourceServer {
	mock := &MockProvisioner_ProvisionResourceServer{ctrl: ctrl}
	mock.recorder = &MockProvisioner_ProvisionResourceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProvisioner_ProvisionResourceServer) EXPECT() *MockProvisioner_ProvisionResourceServerMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockProvisioner_ProvisionResourceServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockProvisioner_ProvisionResourceServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockProvisioner_ProvisionResourceServer)(nil).Context))
}

// RecvMsg mocks base method.
func (m *MockProvisioner_ProvisionResourceServer) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockProvisioner_ProvisionResourceServerMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockProvisioner_ProvisionResourceServer)(nil).RecvMsg), arg0)
}

// Send mocks base method.
func (m *MockProvisioner_ProvisionResourceServer) Send(arg0 *tfplugin5.ProvisionResource_Response) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockProvisioner_ProvisionResourceServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockProvisioner_ProvisionResourceServer)(nil).Send), arg0)
}

// SendHeader mocks base method.
func (m *MockProvisioner_ProvisionResourceServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockProvisioner_ProvisionResourceServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockProvisioner_ProvisionResourceServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m *MockProvisioner_ProvisionResourceServer) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockProvisioner_ProvisionResourceServerMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockProvisioner_ProvisionResourceServer)(nil).SendMsg), arg0)
}

// SetHeader mocks base method.
func (m *MockProvisioner_ProvisionResourceServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockProvisioner_ProvisionResourceServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockProvisioner_ProvisionResourceServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockProvisioner_ProvisionResourceServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockProvisioner_ProvisionResourceServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockProvisioner_ProvisionResourceServer)(nil).SetTrailer), arg0)
}
