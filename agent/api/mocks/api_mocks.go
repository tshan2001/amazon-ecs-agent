// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.
//

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aws/amazon-ecs-agent/agent/api (interfaces: ECSSDK,ECSSubmitStateSDK,ECSClient,AppnetClient,ECSTaskProtectionSDK)

// Package mock_api is a generated GoMock package.
package mock_api

import (
	api "github.com/aws/amazon-ecs-agent/agent/api"
	serviceconnect "github.com/aws/amazon-ecs-agent/agent/api/serviceconnect"
	ecs "github.com/aws/amazon-ecs-agent/agent/ecs_client/model/ecs"
	gomock "github.com/golang/mock/gomock"
	go0 "github.com/prometheus/client_model/go"
	reflect "reflect"
)

// MockECSSDK is a mock of ECSSDK interface
type MockECSSDK struct {
	ctrl     *gomock.Controller
	recorder *MockECSSDKMockRecorder
}

// MockECSSDKMockRecorder is the mock recorder for MockECSSDK
type MockECSSDKMockRecorder struct {
	mock *MockECSSDK
}

// NewMockECSSDK creates a new mock instance
func NewMockECSSDK(ctrl *gomock.Controller) *MockECSSDK {
	mock := &MockECSSDK{ctrl: ctrl}
	mock.recorder = &MockECSSDKMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockECSSDK) EXPECT() *MockECSSDKMockRecorder {
	return m.recorder
}

// CreateCluster mocks base method
func (m *MockECSSDK) CreateCluster(arg0 *ecs.CreateClusterInput) (*ecs.CreateClusterOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCluster", arg0)
	ret0, _ := ret[0].(*ecs.CreateClusterOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCluster indicates an expected call of CreateCluster
func (mr *MockECSSDKMockRecorder) CreateCluster(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCluster", reflect.TypeOf((*MockECSSDK)(nil).CreateCluster), arg0)
}

// DiscoverPollEndpoint mocks base method
func (m *MockECSSDK) DiscoverPollEndpoint(arg0 *ecs.DiscoverPollEndpointInput) (*ecs.DiscoverPollEndpointOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DiscoverPollEndpoint", arg0)
	ret0, _ := ret[0].(*ecs.DiscoverPollEndpointOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DiscoverPollEndpoint indicates an expected call of DiscoverPollEndpoint
func (mr *MockECSSDKMockRecorder) DiscoverPollEndpoint(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DiscoverPollEndpoint", reflect.TypeOf((*MockECSSDK)(nil).DiscoverPollEndpoint), arg0)
}

// ListTagsForResource mocks base method
func (m *MockECSSDK) ListTagsForResource(arg0 *ecs.ListTagsForResourceInput) (*ecs.ListTagsForResourceOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTagsForResource", arg0)
	ret0, _ := ret[0].(*ecs.ListTagsForResourceOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTagsForResource indicates an expected call of ListTagsForResource
func (mr *MockECSSDKMockRecorder) ListTagsForResource(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTagsForResource", reflect.TypeOf((*MockECSSDK)(nil).ListTagsForResource), arg0)
}

// RegisterContainerInstance mocks base method
func (m *MockECSSDK) RegisterContainerInstance(arg0 *ecs.RegisterContainerInstanceInput) (*ecs.RegisterContainerInstanceOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterContainerInstance", arg0)
	ret0, _ := ret[0].(*ecs.RegisterContainerInstanceOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterContainerInstance indicates an expected call of RegisterContainerInstance
func (mr *MockECSSDKMockRecorder) RegisterContainerInstance(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterContainerInstance", reflect.TypeOf((*MockECSSDK)(nil).RegisterContainerInstance), arg0)
}

// UpdateContainerInstancesState mocks base method
func (m *MockECSSDK) UpdateContainerInstancesState(arg0 *ecs.UpdateContainerInstancesStateInput) (*ecs.UpdateContainerInstancesStateOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateContainerInstancesState", arg0)
	ret0, _ := ret[0].(*ecs.UpdateContainerInstancesStateOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateContainerInstancesState indicates an expected call of UpdateContainerInstancesState
func (mr *MockECSSDKMockRecorder) UpdateContainerInstancesState(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContainerInstancesState", reflect.TypeOf((*MockECSSDK)(nil).UpdateContainerInstancesState), arg0)
}

// MockECSSubmitStateSDK is a mock of ECSSubmitStateSDK interface
type MockECSSubmitStateSDK struct {
	ctrl     *gomock.Controller
	recorder *MockECSSubmitStateSDKMockRecorder
}

// MockECSSubmitStateSDKMockRecorder is the mock recorder for MockECSSubmitStateSDK
type MockECSSubmitStateSDKMockRecorder struct {
	mock *MockECSSubmitStateSDK
}

// NewMockECSSubmitStateSDK creates a new mock instance
func NewMockECSSubmitStateSDK(ctrl *gomock.Controller) *MockECSSubmitStateSDK {
	mock := &MockECSSubmitStateSDK{ctrl: ctrl}
	mock.recorder = &MockECSSubmitStateSDKMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockECSSubmitStateSDK) EXPECT() *MockECSSubmitStateSDKMockRecorder {
	return m.recorder
}

// SubmitAttachmentStateChanges mocks base method
func (m *MockECSSubmitStateSDK) SubmitAttachmentStateChanges(arg0 *ecs.SubmitAttachmentStateChangesInput) (*ecs.SubmitAttachmentStateChangesOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitAttachmentStateChanges", arg0)
	ret0, _ := ret[0].(*ecs.SubmitAttachmentStateChangesOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitAttachmentStateChanges indicates an expected call of SubmitAttachmentStateChanges
func (mr *MockECSSubmitStateSDKMockRecorder) SubmitAttachmentStateChanges(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitAttachmentStateChanges", reflect.TypeOf((*MockECSSubmitStateSDK)(nil).SubmitAttachmentStateChanges), arg0)
}

// SubmitContainerStateChange mocks base method
func (m *MockECSSubmitStateSDK) SubmitContainerStateChange(arg0 *ecs.SubmitContainerStateChangeInput) (*ecs.SubmitContainerStateChangeOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitContainerStateChange", arg0)
	ret0, _ := ret[0].(*ecs.SubmitContainerStateChangeOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitContainerStateChange indicates an expected call of SubmitContainerStateChange
func (mr *MockECSSubmitStateSDKMockRecorder) SubmitContainerStateChange(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitContainerStateChange", reflect.TypeOf((*MockECSSubmitStateSDK)(nil).SubmitContainerStateChange), arg0)
}

// SubmitTaskStateChange mocks base method
func (m *MockECSSubmitStateSDK) SubmitTaskStateChange(arg0 *ecs.SubmitTaskStateChangeInput) (*ecs.SubmitTaskStateChangeOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitTaskStateChange", arg0)
	ret0, _ := ret[0].(*ecs.SubmitTaskStateChangeOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitTaskStateChange indicates an expected call of SubmitTaskStateChange
func (mr *MockECSSubmitStateSDKMockRecorder) SubmitTaskStateChange(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitTaskStateChange", reflect.TypeOf((*MockECSSubmitStateSDK)(nil).SubmitTaskStateChange), arg0)
}

// MockECSClient is a mock of ECSClient interface
type MockECSClient struct {
	ctrl     *gomock.Controller
	recorder *MockECSClientMockRecorder
}

// MockECSClientMockRecorder is the mock recorder for MockECSClient
type MockECSClientMockRecorder struct {
	mock *MockECSClient
}

// NewMockECSClient creates a new mock instance
func NewMockECSClient(ctrl *gomock.Controller) *MockECSClient {
	mock := &MockECSClient{ctrl: ctrl}
	mock.recorder = &MockECSClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockECSClient) EXPECT() *MockECSClientMockRecorder {
	return m.recorder
}

// DiscoverPollEndpoint mocks base method
func (m *MockECSClient) DiscoverPollEndpoint(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DiscoverPollEndpoint", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DiscoverPollEndpoint indicates an expected call of DiscoverPollEndpoint
func (mr *MockECSClientMockRecorder) DiscoverPollEndpoint(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DiscoverPollEndpoint", reflect.TypeOf((*MockECSClient)(nil).DiscoverPollEndpoint), arg0)
}

// DiscoverServiceConnectEndpoint mocks base method
func (m *MockECSClient) DiscoverServiceConnectEndpoint(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DiscoverServiceConnectEndpoint", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DiscoverServiceConnectEndpoint indicates an expected call of DiscoverServiceConnectEndpoint
func (mr *MockECSClientMockRecorder) DiscoverServiceConnectEndpoint(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DiscoverServiceConnectEndpoint", reflect.TypeOf((*MockECSClient)(nil).DiscoverServiceConnectEndpoint), arg0)
}

// DiscoverTelemetryEndpoint mocks base method
func (m *MockECSClient) DiscoverTelemetryEndpoint(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DiscoverTelemetryEndpoint", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DiscoverTelemetryEndpoint indicates an expected call of DiscoverTelemetryEndpoint
func (mr *MockECSClientMockRecorder) DiscoverTelemetryEndpoint(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DiscoverTelemetryEndpoint", reflect.TypeOf((*MockECSClient)(nil).DiscoverTelemetryEndpoint), arg0)
}

// GetResourceTags mocks base method
func (m *MockECSClient) GetResourceTags(arg0 string) ([]*ecs.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResourceTags", arg0)
	ret0, _ := ret[0].([]*ecs.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResourceTags indicates an expected call of GetResourceTags
func (mr *MockECSClientMockRecorder) GetResourceTags(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResourceTags", reflect.TypeOf((*MockECSClient)(nil).GetResourceTags), arg0)
}

// RegisterContainerInstance mocks base method
func (m *MockECSClient) RegisterContainerInstance(arg0 string, arg1 []*ecs.Attribute, arg2 []*ecs.Tag, arg3 string, arg4 []*ecs.PlatformDevice, arg5 string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterContainerInstance", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RegisterContainerInstance indicates an expected call of RegisterContainerInstance
func (mr *MockECSClientMockRecorder) RegisterContainerInstance(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterContainerInstance", reflect.TypeOf((*MockECSClient)(nil).RegisterContainerInstance), arg0, arg1, arg2, arg3, arg4, arg5)
}

// SubmitAttachmentStateChange mocks base method
func (m *MockECSClient) SubmitAttachmentStateChange(arg0 api.AttachmentStateChange) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitAttachmentStateChange", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SubmitAttachmentStateChange indicates an expected call of SubmitAttachmentStateChange
func (mr *MockECSClientMockRecorder) SubmitAttachmentStateChange(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitAttachmentStateChange", reflect.TypeOf((*MockECSClient)(nil).SubmitAttachmentStateChange), arg0)
}

// SubmitContainerStateChange mocks base method
func (m *MockECSClient) SubmitContainerStateChange(arg0 api.ContainerStateChange) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitContainerStateChange", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SubmitContainerStateChange indicates an expected call of SubmitContainerStateChange
func (mr *MockECSClientMockRecorder) SubmitContainerStateChange(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitContainerStateChange", reflect.TypeOf((*MockECSClient)(nil).SubmitContainerStateChange), arg0)
}

// SubmitTaskStateChange mocks base method
func (m *MockECSClient) SubmitTaskStateChange(arg0 api.TaskStateChange) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitTaskStateChange", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SubmitTaskStateChange indicates an expected call of SubmitTaskStateChange
func (mr *MockECSClientMockRecorder) SubmitTaskStateChange(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitTaskStateChange", reflect.TypeOf((*MockECSClient)(nil).SubmitTaskStateChange), arg0)
}

// UpdateContainerInstancesState mocks base method
func (m *MockECSClient) UpdateContainerInstancesState(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateContainerInstancesState", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateContainerInstancesState indicates an expected call of UpdateContainerInstancesState
func (mr *MockECSClientMockRecorder) UpdateContainerInstancesState(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContainerInstancesState", reflect.TypeOf((*MockECSClient)(nil).UpdateContainerInstancesState), arg0, arg1)
}

// MockAppnetClient is a mock of AppnetClient interface
type MockAppnetClient struct {
	ctrl     *gomock.Controller
	recorder *MockAppnetClientMockRecorder
}

// MockAppnetClientMockRecorder is the mock recorder for MockAppnetClient
type MockAppnetClientMockRecorder struct {
	mock *MockAppnetClient
}

// NewMockAppnetClient creates a new mock instance
func NewMockAppnetClient(ctrl *gomock.Controller) *MockAppnetClient {
	mock := &MockAppnetClient{ctrl: ctrl}
	mock.recorder = &MockAppnetClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAppnetClient) EXPECT() *MockAppnetClientMockRecorder {
	return m.recorder
}

// DrainInboundConnections mocks base method
func (m *MockAppnetClient) DrainInboundConnections(arg0 serviceconnect.RuntimeConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DrainInboundConnections", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DrainInboundConnections indicates an expected call of DrainInboundConnections
func (mr *MockAppnetClientMockRecorder) DrainInboundConnections(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DrainInboundConnections", reflect.TypeOf((*MockAppnetClient)(nil).DrainInboundConnections), arg0)
}

// GetStats mocks base method
func (m *MockAppnetClient) GetStats(arg0 serviceconnect.RuntimeConfig) (map[string]*go0.MetricFamily, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStats", arg0)
	ret0, _ := ret[0].(map[string]*go0.MetricFamily)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStats indicates an expected call of GetStats
func (mr *MockAppnetClientMockRecorder) GetStats(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStats", reflect.TypeOf((*MockAppnetClient)(nil).GetStats), arg0)
}

// MockECSTaskProtectionSDK is a mock of ECSTaskProtectionSDK interface
type MockECSTaskProtectionSDK struct {
	ctrl     *gomock.Controller
	recorder *MockECSTaskProtectionSDKMockRecorder
}

// MockECSTaskProtectionSDKMockRecorder is the mock recorder for MockECSTaskProtectionSDK
type MockECSTaskProtectionSDKMockRecorder struct {
	mock *MockECSTaskProtectionSDK
}

// NewMockECSTaskProtectionSDK creates a new mock instance
func NewMockECSTaskProtectionSDK(ctrl *gomock.Controller) *MockECSTaskProtectionSDK {
	mock := &MockECSTaskProtectionSDK{ctrl: ctrl}
	mock.recorder = &MockECSTaskProtectionSDKMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockECSTaskProtectionSDK) EXPECT() *MockECSTaskProtectionSDKMockRecorder {
	return m.recorder
}

// GetTaskProtection mocks base method
func (m *MockECSTaskProtectionSDK) GetTaskProtection(arg0 *ecs.GetTaskProtectionInput) (*ecs.GetTaskProtectionOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTaskProtection", arg0)
	ret0, _ := ret[0].(*ecs.GetTaskProtectionOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskProtection indicates an expected call of GetTaskProtection
func (mr *MockECSTaskProtectionSDKMockRecorder) GetTaskProtection(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskProtection", reflect.TypeOf((*MockECSTaskProtectionSDK)(nil).GetTaskProtection), arg0)
}

// UpdateTaskProtection mocks base method
func (m *MockECSTaskProtectionSDK) UpdateTaskProtection(arg0 *ecs.UpdateTaskProtectionInput) (*ecs.UpdateTaskProtectionOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTaskProtection", arg0)
	ret0, _ := ret[0].(*ecs.UpdateTaskProtectionOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTaskProtection indicates an expected call of UpdateTaskProtection
func (mr *MockECSTaskProtectionSDKMockRecorder) UpdateTaskProtection(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTaskProtection", reflect.TypeOf((*MockECSTaskProtectionSDK)(nil).UpdateTaskProtection), arg0)
}
