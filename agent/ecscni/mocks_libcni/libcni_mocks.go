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
// Source: github.com/containernetworking/cni/libcni (interfaces: CNI)

// Package mock_libcni is a generated GoMock package.
package mock_libcni

import (
	context "context"
	reflect "reflect"

	libcni "github.com/containernetworking/cni/libcni"
	types "github.com/containernetworking/cni/pkg/types"
	gomock "github.com/golang/mock/gomock"
)

// MockCNI is a mock of CNI interface
type MockCNI struct {
	ctrl     *gomock.Controller
	recorder *MockCNIMockRecorder
}

// MockCNIMockRecorder is the mock recorder for MockCNI
type MockCNIMockRecorder struct {
	mock *MockCNI
}

// NewMockCNI creates a new mock instance
func NewMockCNI(ctrl *gomock.Controller) *MockCNI {
	mock := &MockCNI{ctrl: ctrl}
	mock.recorder = &MockCNIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCNI) EXPECT() *MockCNIMockRecorder {
	return m.recorder
}

// AddNetwork mocks base method
func (m *MockCNI) AddNetwork(arg0 context.Context, arg1 *libcni.NetworkConfig, arg2 *libcni.RuntimeConf) (types.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNetwork", arg0, arg1, arg2)
	ret0, _ := ret[0].(types.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddNetwork indicates an expected call of AddNetwork
func (mr *MockCNIMockRecorder) AddNetwork(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNetwork", reflect.TypeOf((*MockCNI)(nil).AddNetwork), arg0, arg1, arg2)
}

// AddNetworkList mocks base method
func (m *MockCNI) AddNetworkList(arg0 context.Context, arg1 *libcni.NetworkConfigList, arg2 *libcni.RuntimeConf) (types.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNetworkList", arg0, arg1, arg2)
	ret0, _ := ret[0].(types.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddNetworkList indicates an expected call of AddNetworkList
func (mr *MockCNIMockRecorder) AddNetworkList(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNetworkList", reflect.TypeOf((*MockCNI)(nil).AddNetworkList), arg0, arg1, arg2)
}

// CheckNetwork mocks base method
func (m *MockCNI) CheckNetwork(arg0 context.Context, arg1 *libcni.NetworkConfig, arg2 *libcni.RuntimeConf) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckNetwork", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckNetwork indicates an expected call of CheckNetwork
func (mr *MockCNIMockRecorder) CheckNetwork(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckNetwork", reflect.TypeOf((*MockCNI)(nil).CheckNetwork), arg0, arg1, arg2)
}

// CheckNetworkList mocks base method
func (m *MockCNI) CheckNetworkList(arg0 context.Context, arg1 *libcni.NetworkConfigList, arg2 *libcni.RuntimeConf) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckNetworkList", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckNetworkList indicates an expected call of CheckNetworkList
func (mr *MockCNIMockRecorder) CheckNetworkList(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckNetworkList", reflect.TypeOf((*MockCNI)(nil).CheckNetworkList), arg0, arg1, arg2)
}

// DelNetwork mocks base method
func (m *MockCNI) DelNetwork(arg0 context.Context, arg1 *libcni.NetworkConfig, arg2 *libcni.RuntimeConf) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DelNetwork", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DelNetwork indicates an expected call of DelNetwork
func (mr *MockCNIMockRecorder) DelNetwork(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DelNetwork", reflect.TypeOf((*MockCNI)(nil).DelNetwork), arg0, arg1, arg2)
}

// DelNetworkList mocks base method
func (m *MockCNI) DelNetworkList(arg0 context.Context, arg1 *libcni.NetworkConfigList, arg2 *libcni.RuntimeConf) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DelNetworkList", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DelNetworkList indicates an expected call of DelNetworkList
func (mr *MockCNIMockRecorder) DelNetworkList(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DelNetworkList", reflect.TypeOf((*MockCNI)(nil).DelNetworkList), arg0, arg1, arg2)
}

// GetNetworkCachedConfig mocks base method
func (m *MockCNI) GetNetworkCachedConfig(arg0 *libcni.NetworkConfig, arg1 *libcni.RuntimeConf) ([]byte, *libcni.RuntimeConf, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetworkCachedConfig", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(*libcni.RuntimeConf)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetNetworkCachedConfig indicates an expected call of GetNetworkCachedConfig
func (mr *MockCNIMockRecorder) GetNetworkCachedConfig(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetworkCachedConfig", reflect.TypeOf((*MockCNI)(nil).GetNetworkCachedConfig), arg0, arg1)
}

// GetNetworkCachedResult mocks base method
func (m *MockCNI) GetNetworkCachedResult(arg0 *libcni.NetworkConfig, arg1 *libcni.RuntimeConf) (types.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetworkCachedResult", arg0, arg1)
	ret0, _ := ret[0].(types.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNetworkCachedResult indicates an expected call of GetNetworkCachedResult
func (mr *MockCNIMockRecorder) GetNetworkCachedResult(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetworkCachedResult", reflect.TypeOf((*MockCNI)(nil).GetNetworkCachedResult), arg0, arg1)
}

// GetNetworkListCachedConfig mocks base method
func (m *MockCNI) GetNetworkListCachedConfig(arg0 *libcni.NetworkConfigList, arg1 *libcni.RuntimeConf) ([]byte, *libcni.RuntimeConf, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetworkListCachedConfig", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(*libcni.RuntimeConf)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetNetworkListCachedConfig indicates an expected call of GetNetworkListCachedConfig
func (mr *MockCNIMockRecorder) GetNetworkListCachedConfig(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetworkListCachedConfig", reflect.TypeOf((*MockCNI)(nil).GetNetworkListCachedConfig), arg0, arg1)
}

// GetNetworkListCachedResult mocks base method
func (m *MockCNI) GetNetworkListCachedResult(arg0 *libcni.NetworkConfigList, arg1 *libcni.RuntimeConf) (types.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetworkListCachedResult", arg0, arg1)
	ret0, _ := ret[0].(types.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNetworkListCachedResult indicates an expected call of GetNetworkListCachedResult
func (mr *MockCNIMockRecorder) GetNetworkListCachedResult(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetworkListCachedResult", reflect.TypeOf((*MockCNI)(nil).GetNetworkListCachedResult), arg0, arg1)
}

// ValidateNetwork mocks base method
func (m *MockCNI) ValidateNetwork(arg0 context.Context, arg1 *libcni.NetworkConfig) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateNetwork", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateNetwork indicates an expected call of ValidateNetwork
func (mr *MockCNIMockRecorder) ValidateNetwork(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateNetwork", reflect.TypeOf((*MockCNI)(nil).ValidateNetwork), arg0, arg1)
}

// ValidateNetworkList mocks base method
func (m *MockCNI) ValidateNetworkList(arg0 context.Context, arg1 *libcni.NetworkConfigList) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateNetworkList", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateNetworkList indicates an expected call of ValidateNetworkList
func (mr *MockCNIMockRecorder) ValidateNetworkList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateNetworkList", reflect.TypeOf((*MockCNI)(nil).ValidateNetworkList), arg0, arg1)
}
