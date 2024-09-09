//go:build unit
// +build unit

// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_metrics "github.com/aws/amazon-ecs-agent/ecs-agent/metrics/mocks"
	"github.com/aws/amazon-ecs-agent/ecs-agent/tmds/handlers/fault/v1/types"
	v2 "github.com/aws/amazon-ecs-agent/ecs-agent/tmds/handlers/v2"
	state "github.com/aws/amazon-ecs-agent/ecs-agent/tmds/handlers/v4/state"
	mock_state "github.com/aws/amazon-ecs-agent/ecs-agent/tmds/handlers/v4/state/mocks"
	mock_execwrapper "github.com/aws/amazon-ecs-agent/ecs-agent/utils/execwrapper/mocks"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	endpointId                 = "endpointId"
	port                       = 1234
	protocol                   = "tcp"
	trafficType                = "ingress"
	delayMilliseconds          = 123456789
	jitterMilliseconds         = 4567
	lossPercent                = 6
	taskARN                    = "taskArn"
	awsvpcNetworkMode          = "awsvpc"
	deviceName                 = "eth0"
	invalidNetworkMode         = "invalid"
	iptablesChainNotFoundError = "iptables: Bad rule (does a matching rule exist in that chain?)."
)

var (
	noDeviceNameInNetworkInterfaces = []*state.NetworkInterface{
		{
			DeviceName: "",
		},
	}

	happyNetworkInterfaces = []*state.NetworkInterface{
		{
			DeviceName: deviceName,
		},
	}

	happyNetworkNamespaces = []*state.NetworkNamespace{
		{
			Path:              "/some/path",
			NetworkInterfaces: happyNetworkInterfaces,
		},
	}

	noPathInNetworkNamespaces = []*state.NetworkNamespace{
		{
			Path:              "",
			NetworkInterfaces: happyNetworkInterfaces,
		},
	}

	happyTaskNetworkConfig = state.TaskNetworkConfig{
		NetworkMode:       awsvpcNetworkMode,
		NetworkNamespaces: happyNetworkNamespaces,
	}

	happyTaskResponse = state.TaskResponse{
		TaskResponse:          &v2.TaskResponse{TaskARN: taskARN},
		TaskNetworkConfig:     &happyTaskNetworkConfig,
		FaultInjectionEnabled: true,
	}

	happyBlackHolePortReqBody = map[string]interface{}{
		"Port":        port,
		"Protocol":    protocol,
		"TrafficType": trafficType,
	}

	ipSources = []string{"52.95.154.1", "52.95.154.2"}
)

type networkFaultInjectionTestCase struct {
	name                      string
	expectedStatusCode        int
	requestBody               interface{}
	expectedResponseBody      types.NetworkFaultInjectionResponse
	setAgentStateExpectations func(agentState *mock_state.MockAgentState)
	setExecExpectations       func(exec *mock_execwrapper.MockExec, ctrl *gomock.Controller)
}

// Tests the path for Fault Network Faults API
func TestFaultBlackholeFaultPath(t *testing.T) {
	assert.Equal(t, "/api/{endpointContainerIDMuxName:[^/]*}/fault/v1/network-blackhole-port", NetworkFaultPath(types.BlackHolePortFaultType))
}

func TestFaultLatencyFaultPath(t *testing.T) {
	assert.Equal(t, "/api/{endpointContainerIDMuxName:[^/]*}/fault/v1/network-latency", NetworkFaultPath(types.LatencyFaultType))
}

func TestFaultPacketLossFaultPath(t *testing.T) {
	assert.Equal(t, "/api/{endpointContainerIDMuxName:[^/]*}/fault/v1/network-packet-loss", NetworkFaultPath(types.PacketLossFaultType))
}

func generateNetworkBlackHolePortTestCases(name string) []networkFaultInjectionTestCase {
	tcs := []networkFaultInjectionTestCase{
		{
			name:                 fmt.Sprintf("%s no request body", name),
			expectedStatusCode:   400,
			requestBody:          nil,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("required request body is missing"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s malformed request body", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"Port":        "incorrect typing",
				"Protocol":    protocol,
				"TrafficType": trafficType,
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("json: cannot unmarshal string into Go struct field NetworkBlackholePortRequest.Port of type uint16"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s incomplete request body", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"Port":     port,
				"Protocol": protocol,
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("required parameter TrafficType is missing"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s empty value request body", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"Port":        port,
				"Protocol":    protocol,
				"TrafficType": "",
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("required parameter TrafficType is missing"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s invalid port value request body", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"Port":        port,
				"Protocol":    "invalid",
				"TrafficType": trafficType,
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("invalid value invalid for parameter Protocol"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s invalid traffic type value request body", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"Port":        port,
				"Protocol":    protocol,
				"TrafficType": "invalid",
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("invalid value invalid for parameter TrafficType"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:                 fmt.Sprintf("%s task lookup fail", name),
			expectedStatusCode:   404,
			requestBody:          happyBlackHolePortReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(fmt.Sprintf("unable to lookup container: %s", endpointId)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(state.TaskResponse{}, state.NewErrorLookupFailure("task lookup failed")).
					Times(1)
			},
		},
		{
			name:                 fmt.Sprintf("%s task metadata fetch fail", name),
			expectedStatusCode:   500,
			requestBody:          happyBlackHolePortReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(fmt.Sprintf("unable to obtain container metadata for container: %s", endpointId)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, state.NewErrorMetadataFetchFailure(
					"Unable to generate metadata for task")).
					Times(1)
			},
		},
		{
			name:                 fmt.Sprintf("%s task metadata unknown fail", name),
			expectedStatusCode:   500,
			requestBody:          happyBlackHolePortReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(fmt.Sprintf("failed to get task metadata due to internal server error for container: %s", endpointId)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(state.TaskResponse{}, errors.New("unknown error")).
					Times(1)
			},
		},
		{
			name:                 fmt.Sprintf("%s fault injection disabled", name),
			expectedStatusCode:   400,
			requestBody:          happyBlackHolePortReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(fmt.Sprintf(faultInjectionEnabledError, taskARN)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{
					TaskResponse:          &v2.TaskResponse{TaskARN: taskARN},
					FaultInjectionEnabled: false,
				}, nil).Times(1)
			},
		},
		{
			name:                 fmt.Sprintf("%s invalid network mode", name),
			expectedStatusCode:   400,
			requestBody:          happyBlackHolePortReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(fmt.Sprintf(invalidNetworkModeError, invalidNetworkMode)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{
					TaskResponse:          &v2.TaskResponse{TaskARN: taskARN},
					FaultInjectionEnabled: true,
					TaskNetworkConfig: &state.TaskNetworkConfig{
						NetworkMode:       invalidNetworkMode,
						NetworkNamespaces: happyNetworkNamespaces,
					},
				}, nil).Times(1)
			},
		},
		{
			name:               fmt.Sprintf("%s empty task network config", name),
			expectedStatusCode: 500,
			requestBody:        happyBlackHolePortReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(
				fmt.Sprintf("failed to get task metadata due to internal server error for container: %s", endpointId)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{
					TaskResponse:          &v2.TaskResponse{TaskARN: taskARN},
					FaultInjectionEnabled: true,
					TaskNetworkConfig:     nil,
				}, nil).Times(1)
			},
		},
		{
			name:               fmt.Sprintf("%s no task network namespace", name),
			expectedStatusCode: 500,
			requestBody:        happyBlackHolePortReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(
				fmt.Sprintf("failed to get task metadata due to internal server error for container: %s", endpointId)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{
					TaskResponse:          &v2.TaskResponse{TaskARN: taskARN},
					FaultInjectionEnabled: true,
					TaskNetworkConfig: &state.TaskNetworkConfig{
						NetworkMode:       awsvpcNetworkMode,
						NetworkNamespaces: nil,
					},
				}, nil).Times(1)
			},
		},
		{
			name:               fmt.Sprintf("%s no path in task network config", name),
			expectedStatusCode: 500,
			requestBody:        happyBlackHolePortReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(
				fmt.Sprintf("failed to get task metadata due to internal server error for container: %s", endpointId)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{
					TaskResponse:          &v2.TaskResponse{TaskARN: taskARN},
					FaultInjectionEnabled: true,
					TaskNetworkConfig: &state.TaskNetworkConfig{
						NetworkMode:       awsvpcNetworkMode,
						NetworkNamespaces: noPathInNetworkNamespaces,
					},
				}, nil).Times(1)
			},
		},
		{
			name:               fmt.Sprintf("%s no device name in task network config", name),
			expectedStatusCode: 500,
			requestBody:        happyBlackHolePortReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(
				fmt.Sprintf("failed to get task metadata due to internal server error for container: %s", endpointId)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{
					TaskResponse:          &v2.TaskResponse{TaskARN: taskARN},
					FaultInjectionEnabled: true,
					TaskNetworkConfig: &state.TaskNetworkConfig{
						NetworkMode: awsvpcNetworkMode,
						NetworkNamespaces: []*state.NetworkNamespace{
							{
								Path:              "/path",
								NetworkInterfaces: noDeviceNameInNetworkInterfaces,
							},
						},
					},
				}, nil).Times(1)
			},
		},
	}
	return tcs
}

func generateStartBlackHolePortFaultTestCases() []networkFaultInjectionTestCase {
	commonTcs := generateNetworkBlackHolePortTestCases("start blackhole port")
	tcs := []networkFaultInjectionTestCase{
		{
			name:                 "start blackhole port success running",
			expectedStatusCode:   200,
			requestBody:          happyBlackHolePortReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionSuccessResponse("running"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(happyTaskResponse, nil).
					Times(1)
			},
			setExecExpectations: func(exec *mock_execwrapper.MockExec, ctrl *gomock.Controller) {
			},
		},
		{
			name:               "start blackhole unknown request body",
			expectedStatusCode: 200,
			requestBody: map[string]interface{}{
				"Port":        port,
				"Protocol":    protocol,
				"TrafficType": trafficType,
				"Unknown":     "",
			},
			expectedResponseBody: types.NewNetworkFaultInjectionSuccessResponse("running"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(happyTaskResponse, nil).
					Times(1)
			},
			setExecExpectations: func(exec *mock_execwrapper.MockExec, ctrl *gomock.Controller) {
			},
		},
	}

	return append(tcs, commonTcs...)
}

func generateStopBlackHolePortFaultTestCases() []networkFaultInjectionTestCase {
	commonTcs := generateNetworkBlackHolePortTestCases("stop blackhole port")
	tcs := []networkFaultInjectionTestCase{
		{
			name:                 "stop blackhole port success running",
			expectedStatusCode:   200,
			requestBody:          happyBlackHolePortReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionSuccessResponse("stopped"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(happyTaskResponse, nil).
					Times(1)
			},
			setExecExpectations: func(exec *mock_execwrapper.MockExec, ctrl *gomock.Controller) {
			},
		},
		{
			name:               "stop blackhole unknown request body",
			expectedStatusCode: 200,
			requestBody: map[string]interface{}{
				"Port":        port,
				"Protocol":    protocol,
				"TrafficType": trafficType,
				"Unknown":     "",
			},
			expectedResponseBody: types.NewNetworkFaultInjectionSuccessResponse("stopped"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(happyTaskResponse, nil).
					Times(1)
			},
			setExecExpectations: func(exec *mock_execwrapper.MockExec, ctrl *gomock.Controller) {
			},
		},
	}
	return append(tcs, commonTcs...)
}

func generateCheckBlackHolePortFaultStatusTestCases() []networkFaultInjectionTestCase {
	commonTcs := generateNetworkBlackHolePortTestCases("check blackhole port")
	tcs := []networkFaultInjectionTestCase{
		{
			name:                 "check blackhole port success running",
			expectedStatusCode:   200,
			requestBody:          happyBlackHolePortReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionSuccessResponse("running"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(happyTaskResponse, nil).
					Times(1)
			},
			setExecExpectations: func(exec *mock_execwrapper.MockExec, ctrl *gomock.Controller) {
				cmdExec := mock_execwrapper.NewMockCmd(ctrl)
				gomock.InOrder(
					exec.EXPECT().CommandContext(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(cmdExec),
					cmdExec.EXPECT().CombinedOutput().Times(1).Return([]byte{}, nil),
				)
			},
		},
		{
			name:               "check blackhole unknown request body",
			expectedStatusCode: 200,
			requestBody: map[string]interface{}{
				"Port":        port,
				"Protocol":    protocol,
				"TrafficType": trafficType,
				"Unknown":     "",
			},
			expectedResponseBody: types.NewNetworkFaultInjectionSuccessResponse("running"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(happyTaskResponse, nil).
					Times(1)
			},
			setExecExpectations: func(exec *mock_execwrapper.MockExec, ctrl *gomock.Controller) {
				cmdExec := mock_execwrapper.NewMockCmd(ctrl)
				gomock.InOrder(
					exec.EXPECT().CommandContext(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(cmdExec),
					cmdExec.EXPECT().CombinedOutput().Times(1).Return([]byte{}, nil),
				)
			},
		},
		{
			name:                 "check blackhole port success not running",
			expectedStatusCode:   200,
			requestBody:          happyBlackHolePortReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionSuccessResponse("not-running"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(happyTaskResponse, nil).
					Times(1)
			},
			setExecExpectations: func(exec *mock_execwrapper.MockExec, ctrl *gomock.Controller) {
				cmdExec := mock_execwrapper.NewMockCmd(ctrl)
				gomock.InOrder(
					exec.EXPECT().CommandContext(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(cmdExec),
					cmdExec.EXPECT().CombinedOutput().Times(1).Return([]byte(iptablesChainNotFoundError), errors.New("exit status 1")),
					exec.EXPECT().ConvertToExitError(gomock.Any()).Times(1).Return(nil, true),
					exec.EXPECT().GetExitCode(gomock.Any()).Times(1).Return(1),
				)
			},
		},
		{
			name:                 "check blackhole port failure",
			expectedStatusCode:   500,
			requestBody:          happyBlackHolePortReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("internal error"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(happyTaskResponse, nil).
					Times(1)
			},
			setExecExpectations: func(exec *mock_execwrapper.MockExec, ctrl *gomock.Controller) {
				cmdExec := mock_execwrapper.NewMockCmd(ctrl)
				gomock.InOrder(
					exec.EXPECT().CommandContext(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(cmdExec),
					cmdExec.EXPECT().CombinedOutput().Times(1).Return([]byte("internal error"), errors.New("exit 2")),
					exec.EXPECT().ConvertToExitError(gomock.Any()).Times(1).Return(nil, false),
				)
			},
		},
	}

	return append(tcs, commonTcs...)
}

func TestStartNetworkBlackHolePort(t *testing.T) {
	tcs := generateStartBlackHolePortFaultTestCases()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// Mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			agentState := mock_state.NewMockAgentState(ctrl)
			metricsFactory := mock_metrics.NewMockEntryFactory(ctrl)
			execWrapper := mock_execwrapper.NewMockExec(ctrl)

			if tc.setAgentStateExpectations != nil {
				tc.setAgentStateExpectations(agentState)
			}

			router := mux.NewRouter()
			handler := New(agentState, metricsFactory, execWrapper)
			router.HandleFunc(
				NetworkFaultPath(types.BlackHolePortFaultType),
				handler.StartNetworkBlackholePort(),
			).Methods("PUT")

			method := "PUT"
			var requestBody io.Reader
			if tc.requestBody != nil {
				reqBodyBytes, err := json.Marshal(tc.requestBody)
				require.NoError(t, err)
				requestBody = bytes.NewReader(reqBodyBytes)
			}
			req, err := http.NewRequest(method, fmt.Sprintf("/api/%s/fault/v1/network-blackhole-port", endpointId),
				requestBody)
			require.NoError(t, err)

			// Send the request and record the response
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			var actualResponseBody types.NetworkFaultInjectionResponse
			err = json.Unmarshal(recorder.Body.Bytes(), &actualResponseBody)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Equal(t, tc.expectedResponseBody, actualResponseBody)

		})
	}
}

func TestStopNetworkBlackHolePort(t *testing.T) {
	tcs := generateStopBlackHolePortFaultTestCases()
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// Mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			agentState := mock_state.NewMockAgentState(ctrl)
			metricsFactory := mock_metrics.NewMockEntryFactory(ctrl)
			execWrapper := mock_execwrapper.NewMockExec(ctrl)

			if tc.setAgentStateExpectations != nil {
				tc.setAgentStateExpectations(agentState)
			}

			router := mux.NewRouter()
			handler := New(agentState, metricsFactory, execWrapper)
			router.HandleFunc(
				NetworkFaultPath(types.BlackHolePortFaultType),
				handler.StopNetworkBlackHolePort(),
			).Methods("DELETE")

			method := "DELETE"
			var requestBody io.Reader
			if tc.requestBody != nil {
				reqBodyBytes, err := json.Marshal(tc.requestBody)
				require.NoError(t, err)
				requestBody = bytes.NewReader(reqBodyBytes)
			}
			req, err := http.NewRequest(method, fmt.Sprintf("/api/%s/fault/v1/network-blackhole-port", endpointId),
				requestBody)
			require.NoError(t, err)

			// Send the request and record the response
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			var actualResponseBody types.NetworkFaultInjectionResponse
			err = json.Unmarshal(recorder.Body.Bytes(), &actualResponseBody)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Equal(t, tc.expectedResponseBody, actualResponseBody)

		})
	}
}

func TestCheckNetworkBlackHolePort(t *testing.T) {
	tcs := generateCheckBlackHolePortFaultStatusTestCases()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// Mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			agentState := mock_state.NewMockAgentState(ctrl)
			metricsFactory := mock_metrics.NewMockEntryFactory(ctrl)
			execWrapper := mock_execwrapper.NewMockExec(ctrl)

			if tc.setAgentStateExpectations != nil {
				tc.setAgentStateExpectations(agentState)
			}

			if tc.setExecExpectations != nil {
				tc.setExecExpectations(execWrapper, ctrl)
			}

			router := mux.NewRouter()
			handler := New(agentState, metricsFactory, execWrapper)
			router.HandleFunc(
				NetworkFaultPath(types.BlackHolePortFaultType),
				handler.CheckNetworkBlackHolePort(),
			).Methods("GET")

			method := "GET"
			var requestBody io.Reader
			if tc.requestBody != nil {
				reqBodyBytes, err := json.Marshal(tc.requestBody)
				require.NoError(t, err)
				requestBody = bytes.NewReader(reqBodyBytes)
			}
			req, err := http.NewRequest(method, fmt.Sprintf("/api/%s/fault/v1/network-blackhole-port", endpointId),
				requestBody)
			require.NoError(t, err)

			// Send the request and record the response
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			var actualResponseBody types.NetworkFaultInjectionResponse
			err = json.Unmarshal(recorder.Body.Bytes(), &actualResponseBody)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Equal(t, tc.expectedResponseBody, actualResponseBody)

		})
	}
}

func generateNetworkLatencyTestCases(name, expectedHappyResponseBody string) []networkFaultInjectionTestCase {
	happyNetworkLatencyReqBody := map[string]interface{}{
		"DelayMilliseconds":  delayMilliseconds,
		"JitterMilliseconds": jitterMilliseconds,
		"Sources":            ipSources,
	}
	tcs := []networkFaultInjectionTestCase{
		{
			name:                 fmt.Sprintf("%s success", name),
			expectedStatusCode:   200,
			requestBody:          happyNetworkLatencyReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionSuccessResponse(expectedHappyResponseBody),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(happyTaskResponse, nil).
					Times(1)
			},
		},
		{
			name:               fmt.Sprintf("%s unknown request body", name),
			expectedStatusCode: 200,
			requestBody: map[string]interface{}{
				"DelayMilliseconds":  delayMilliseconds,
				"JitterMilliseconds": jitterMilliseconds,
				"Sources":            ipSources,
				"Unknown":            "",
			},
			expectedResponseBody: types.NewNetworkFaultInjectionSuccessResponse(expectedHappyResponseBody),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(happyTaskResponse, nil).
					Times(1)
			},
		},
		{
			name:                 fmt.Sprintf("%s no request body", name),
			expectedStatusCode:   400,
			requestBody:          nil,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("required request body is missing"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s malformed request body 1", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"DelayMilliseconds":  "incorrect-field",
				"JitterMilliseconds": jitterMilliseconds,
				"Sources":            ipSources,
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("json: cannot unmarshal string into Go struct field NetworkLatencyRequest.DelayMilliseconds of type uint64"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(state.TaskResponse{}, nil).
					Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s malformed request body 2", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"DelayMilliseconds":  delayMilliseconds,
				"JitterMilliseconds": "incorrect-field",
				"Sources":            ipSources,
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("json: cannot unmarshal string into Go struct field NetworkLatencyRequest.JitterMilliseconds of type uint64"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(state.TaskResponse{}, nil).
					Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s malformed request body 3", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"DelayMilliseconds":  delayMilliseconds,
				"JitterMilliseconds": jitterMilliseconds,
				"Sources":            "",
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("json: cannot unmarshal string into Go struct field NetworkLatencyRequest.Sources of type []*string"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(state.TaskResponse{}, nil).
					Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s incomplete request body 1", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"DelayMilliseconds":  delayMilliseconds,
				"JitterMilliseconds": jitterMilliseconds,
				"Sources":            []string{},
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("required parameter Sources is missing"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(state.TaskResponse{}, nil).
					Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s incomplete request body 2", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"DelayMilliseconds":  delayMilliseconds,
				"JitterMilliseconds": jitterMilliseconds,
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("required parameter Sources is missing"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(state.TaskResponse{}, nil).
					Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s incomplete request body 3", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"DelayMilliseconds": delayMilliseconds,
				"Sources":           ipSources,
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("required parameter JitterMilliseconds is missing"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(state.TaskResponse{}, nil).
					Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s incomplete request body 4", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"JitterMilliseconds": jitterMilliseconds,
				"Sources":            ipSources,
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("required parameter DelayMilliseconds is missing"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(state.TaskResponse{}, nil).
					Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s invalid IP value in the request body 1", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"DelayMilliseconds":  delayMilliseconds,
				"JitterMilliseconds": jitterMilliseconds,
				"Sources":            []string{"10.1.2.3.4"},
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("invalid value 10.1.2.3.4 for parameter Sources"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(state.TaskResponse{}, nil).
					Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s invalid IP CIDR block value in the request body 2", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"DelayMilliseconds":  delayMilliseconds,
				"JitterMilliseconds": jitterMilliseconds,
				"Sources":            []string{"52.95.154.0/33"},
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("invalid value 52.95.154.0/33 for parameter Sources"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(state.TaskResponse{}, nil).
					Times(0)
			},
		},
		{
			name:                 fmt.Sprintf("%s task lookup fail", name),
			expectedStatusCode:   404,
			requestBody:          happyNetworkLatencyReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(fmt.Sprintf("unable to lookup container: %s", endpointId)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(state.TaskResponse{}, state.NewErrorLookupFailure("task lookup failed")).
					Times(1)
			},
		},
		{
			name:                 fmt.Sprintf("%s task metadata fetch fail", name),
			expectedStatusCode:   500,
			requestBody:          happyNetworkLatencyReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(fmt.Sprintf("unable to obtain container metadata for container: %s", endpointId)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(state.TaskResponse{}, state.NewErrorMetadataFetchFailure(
						"Unable to generate metadata for task")).
					Times(1)
			},
		},
		{
			name:                 fmt.Sprintf("%s task metadata unknown fail", name),
			expectedStatusCode:   500,
			requestBody:          happyNetworkLatencyReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(fmt.Sprintf("failed to get task metadata due to internal server error for container: %s", endpointId)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).
					Return(state.TaskResponse{}, errors.New("unknown error")).
					Times(1)
			},
		},
		{
			name:                 fmt.Sprintf("%s fault injection disabled", name),
			expectedStatusCode:   400,
			requestBody:          happyNetworkLatencyReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(fmt.Sprintf(faultInjectionEnabledError, taskARN)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{
					TaskResponse:          &v2.TaskResponse{TaskARN: taskARN},
					FaultInjectionEnabled: false,
				}, nil)
			},
		},
		{
			name:                 fmt.Sprintf("%s invalid network mode", name),
			expectedStatusCode:   400,
			requestBody:          happyNetworkLatencyReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(fmt.Sprintf(invalidNetworkModeError, invalidNetworkMode)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{
					TaskResponse:          &v2.TaskResponse{TaskARN: taskARN},
					FaultInjectionEnabled: true,
					TaskNetworkConfig: &state.TaskNetworkConfig{
						NetworkMode:       invalidNetworkMode,
						NetworkNamespaces: happyNetworkNamespaces,
					},
				}, nil)
			},
		},
		{
			name:                 fmt.Sprintf("%s empty task network config", name),
			expectedStatusCode:   500,
			requestBody:          happyNetworkLatencyReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(fmt.Sprintf("failed to get task metadata due to internal server error for container: %s", endpointId)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{
					TaskResponse:          &v2.TaskResponse{TaskARN: taskARN},
					FaultInjectionEnabled: true,
					TaskNetworkConfig:     nil,
				}, nil)
			},
		},
	}
	return tcs
}

func TestStartNetworkLatency(t *testing.T) {
	tcs := generateNetworkLatencyTestCases("start network latency", "running")

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// Mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			agentState := mock_state.NewMockAgentState(ctrl)
			metricsFactory := mock_metrics.NewMockEntryFactory(ctrl)
			execWrapper := mock_execwrapper.NewMockExec(ctrl)

			if tc.setAgentStateExpectations != nil {
				tc.setAgentStateExpectations(agentState)
			}

			router := mux.NewRouter()
			handler := New(agentState, metricsFactory, execWrapper)
			router.HandleFunc(
				NetworkFaultPath(types.LatencyFaultType),
				handler.StartNetworkLatency(),
			).Methods("PUT")

			method := "PUT"
			var requestBody io.Reader
			if tc.requestBody != nil {
				reqBodyBytes, err := json.Marshal(tc.requestBody)
				require.NoError(t, err)
				requestBody = bytes.NewReader(reqBodyBytes)
			}
			req, err := http.NewRequest(method, fmt.Sprintf("/api/%s/fault/v1/network-latency", endpointId),
				requestBody)
			require.NoError(t, err)

			// Send the request and record the response
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			var actualResponseBody types.NetworkFaultInjectionResponse
			err = json.Unmarshal(recorder.Body.Bytes(), &actualResponseBody)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Equal(t, tc.expectedResponseBody, actualResponseBody)

		})
	}
}

func TestStopNetworkLatency(t *testing.T) {
	tcs := generateNetworkLatencyTestCases("stop network latency", "stopped")
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// Mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			agentState := mock_state.NewMockAgentState(ctrl)
			metricsFactory := mock_metrics.NewMockEntryFactory(ctrl)
			execWrapper := mock_execwrapper.NewMockExec(ctrl)

			if tc.setAgentStateExpectations != nil {
				tc.setAgentStateExpectations(agentState)
			}

			router := mux.NewRouter()
			handler := New(agentState, metricsFactory, execWrapper)
			router.HandleFunc(
				NetworkFaultPath(types.LatencyFaultType),
				handler.StopNetworkLatency(),
			).Methods("DELETE")

			method := "DELETE"
			var requestBody io.Reader
			if tc.requestBody != nil {
				reqBodyBytes, err := json.Marshal(tc.requestBody)
				require.NoError(t, err)
				requestBody = bytes.NewReader(reqBodyBytes)
			}
			req, err := http.NewRequest(method, fmt.Sprintf("/api/%s/fault/v1/network-latency", endpointId),
				requestBody)
			require.NoError(t, err)

			// Send the request and record the response
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			var actualResponseBody types.NetworkFaultInjectionResponse
			err = json.Unmarshal(recorder.Body.Bytes(), &actualResponseBody)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Equal(t, tc.expectedResponseBody, actualResponseBody)

		})
	}
}

func TestCheckNetworkLatency(t *testing.T) {
	tcs := generateNetworkLatencyTestCases("check network latency", "running")

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// Mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			agentState := mock_state.NewMockAgentState(ctrl)
			metricsFactory := mock_metrics.NewMockEntryFactory(ctrl)
			execWrapper := mock_execwrapper.NewMockExec(ctrl)

			router := mux.NewRouter()
			handler := New(agentState, metricsFactory, execWrapper)
			if tc.setAgentStateExpectations != nil {
				tc.setAgentStateExpectations(agentState)
			}

			router.HandleFunc(
				NetworkFaultPath(types.LatencyFaultType),
				handler.CheckNetworkLatency(),
			).Methods("GET")

			method := "GET"
			var requestBody io.Reader
			if tc.requestBody != nil {
				reqBodyBytes, err := json.Marshal(tc.requestBody)
				require.NoError(t, err)
				requestBody = bytes.NewReader(reqBodyBytes)
			}
			req, err := http.NewRequest(method, fmt.Sprintf("/api/%s/fault/v1/network-latency", endpointId),
				requestBody)
			require.NoError(t, err)

			// Send the request and record the response
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			var actualResponseBody types.NetworkFaultInjectionResponse
			err = json.Unmarshal(recorder.Body.Bytes(), &actualResponseBody)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Equal(t, tc.expectedResponseBody, actualResponseBody)

		})
	}
}

func generateNetworkPacketLossTestCases(name, expectedHappyResponseBody string) []networkFaultInjectionTestCase {
	happyNetworkPacketLossReqBody := map[string]interface{}{
		"LossPercent": lossPercent,
		"Sources":     ipSources,
	}
	tcs := []networkFaultInjectionTestCase{
		{
			name:                 fmt.Sprintf("%s success", name),
			expectedStatusCode:   200,
			requestBody:          happyNetworkPacketLossReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionSuccessResponse(expectedHappyResponseBody),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(happyTaskResponse, nil)
			},
		},
		{
			name:               fmt.Sprintf("%s unknown request body", name),
			expectedStatusCode: 200,
			requestBody: map[string]interface{}{
				"LossPercent": lossPercent,
				"Sources":     ipSources,
				"Unknown":     "",
			},
			expectedResponseBody: types.NewNetworkFaultInjectionSuccessResponse(expectedHappyResponseBody),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(happyTaskResponse, nil)
			},
		},
		{
			name:                 fmt.Sprintf("%s no request body", name),
			expectedStatusCode:   400,
			requestBody:          nil,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("required request body is missing"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s malformed request body 1", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"LossPercent": "incorrect-field",
				"Sources":     ipSources,
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("json: cannot unmarshal string into Go struct field NetworkPacketLossRequest.LossPercent of type uint64"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s malformed request body 2", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"LossPercent": lossPercent,
				"Sources":     "",
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("json: cannot unmarshal string into Go struct field NetworkPacketLossRequest.Sources of type []*string"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s incomplete request body 1", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"LossPercent": lossPercent,
				"Sources":     []string{},
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("required parameter Sources is missing"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s incomplete request body 2", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"LossPercent": lossPercent,
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("required parameter Sources is missing"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s incomplete request body 3", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"Sources": ipSources,
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("required parameter LossPercent is missing"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s invalid LossPercent in the request body 1", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"LossPercent": -1,
				"Sources":     ipSources,
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("json: cannot unmarshal number -1 into Go struct field NetworkPacketLossRequest.LossPercent of type uint64"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s invalid LossPercent in the request body 2", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"LossPercent": 101,
				"Sources":     ipSources,
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("invalid value 101 for parameter LossPercent"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s invalid LossPercent in the request body 3", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"LossPercent": 0,
				"Sources":     ipSources,
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("invalid value 0 for parameter LossPercent"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s invalid IP value in the request body 1", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"LossPercent": lossPercent,
				"Sources":     []string{"10.1.2.3.4"},
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("invalid value 10.1.2.3.4 for parameter Sources"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:               fmt.Sprintf("%s invalid IP CIDR block value in the request body 2", name),
			expectedStatusCode: 400,
			requestBody: map[string]interface{}{
				"LossPercent": lossPercent,
				"Sources":     []string{"52.95.154.0/33"},
			},
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse("invalid value 52.95.154.0/33 for parameter Sources"),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, nil).Times(0)
			},
		},
		{
			name:                 fmt.Sprintf("%s task lookup fail", name),
			expectedStatusCode:   404,
			requestBody:          happyNetworkPacketLossReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(fmt.Sprintf("unable to lookup container: %s", endpointId)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, state.NewErrorLookupFailure("task lookup failed"))
			},
		},
		{
			name:                 fmt.Sprintf("%s task metadata fetch fail", name),
			expectedStatusCode:   500,
			requestBody:          happyNetworkPacketLossReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(fmt.Sprintf("unable to obtain container metadata for container: %s", endpointId)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, state.NewErrorMetadataFetchFailure(
					"Unable to generate metadata for task"))
			},
		},
		{
			name:                 fmt.Sprintf("%s task metadata unknown fail", name),
			expectedStatusCode:   500,
			requestBody:          happyNetworkPacketLossReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(fmt.Sprintf("failed to get task metadata due to internal server error for container: %s", endpointId)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{}, errors.New("unknown error"))
			},
		},
		{
			name:                 fmt.Sprintf("%s fault injection disabled", name),
			expectedStatusCode:   400,
			requestBody:          happyNetworkPacketLossReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(fmt.Sprintf(faultInjectionEnabledError, taskARN)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{
					TaskResponse:          &v2.TaskResponse{TaskARN: taskARN},
					FaultInjectionEnabled: false,
				}, nil)
			},
		},
		{
			name:                 fmt.Sprintf("%s invalid network mode", name),
			expectedStatusCode:   400,
			requestBody:          happyNetworkPacketLossReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(fmt.Sprintf(invalidNetworkModeError, invalidNetworkMode)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{
					TaskResponse:          &v2.TaskResponse{TaskARN: taskARN},
					FaultInjectionEnabled: true,
					TaskNetworkConfig: &state.TaskNetworkConfig{
						NetworkMode:       invalidNetworkMode,
						NetworkNamespaces: happyNetworkNamespaces,
					},
				}, nil)
			},
		},
		{
			name:                 fmt.Sprintf("%s empty task network config", name),
			expectedStatusCode:   500,
			requestBody:          happyNetworkPacketLossReqBody,
			expectedResponseBody: types.NewNetworkFaultInjectionErrorResponse(fmt.Sprintf("failed to get task metadata due to internal server error for container: %s", endpointId)),
			setAgentStateExpectations: func(agentState *mock_state.MockAgentState) {
				agentState.EXPECT().GetTaskMetadata(endpointId).Return(state.TaskResponse{
					TaskResponse:          &v2.TaskResponse{TaskARN: taskARN},
					FaultInjectionEnabled: true,
					TaskNetworkConfig:     nil,
				}, nil)
			},
		},
	}
	return tcs
}

func TestStartNetworkPacketLoss(t *testing.T) {
	tcs := generateNetworkPacketLossTestCases("start network packet loss", "running")

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// Mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			agentState := mock_state.NewMockAgentState(ctrl)
			metricsFactory := mock_metrics.NewMockEntryFactory(ctrl)
			execWrapper := mock_execwrapper.NewMockExec(ctrl)

			if tc.setAgentStateExpectations != nil {
				tc.setAgentStateExpectations(agentState)
			}

			router := mux.NewRouter()
			handler := New(agentState, metricsFactory, execWrapper)
			router.HandleFunc(
				NetworkFaultPath(types.PacketLossFaultType),
				handler.StartNetworkPacketLoss(),
			).Methods("PUT")

			method := "PUT"
			var requestBody io.Reader
			if tc.requestBody != nil {
				reqBodyBytes, err := json.Marshal(tc.requestBody)
				require.NoError(t, err)
				requestBody = bytes.NewReader(reqBodyBytes)
			}
			req, err := http.NewRequest(method, fmt.Sprintf("/api/%s/fault/v1/network-packet-loss", endpointId),
				requestBody)
			require.NoError(t, err)

			// Send the request and record the response
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			var actualResponseBody types.NetworkFaultInjectionResponse
			err = json.Unmarshal(recorder.Body.Bytes(), &actualResponseBody)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Equal(t, tc.expectedResponseBody, actualResponseBody)

		})
	}
}

func TestStopNetworkPacketLoss(t *testing.T) {
	tcs := generateNetworkPacketLossTestCases("stop network packet loss", "stopped")
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// Mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			agentState := mock_state.NewMockAgentState(ctrl)
			metricsFactory := mock_metrics.NewMockEntryFactory(ctrl)
			execWrapper := mock_execwrapper.NewMockExec(ctrl)

			if tc.setAgentStateExpectations != nil {
				tc.setAgentStateExpectations(agentState)
			}

			router := mux.NewRouter()
			handler := New(agentState, metricsFactory, execWrapper)
			router.HandleFunc(
				NetworkFaultPath(types.PacketLossFaultType),
				handler.StopNetworkPacketLoss(),
			).Methods("DELETE")

			method := "DELETE"
			var requestBody io.Reader
			if tc.requestBody != nil {
				reqBodyBytes, err := json.Marshal(tc.requestBody)
				require.NoError(t, err)
				requestBody = bytes.NewReader(reqBodyBytes)
			}
			req, err := http.NewRequest(method, fmt.Sprintf("/api/%s/fault/v1/network-packet-loss", endpointId),
				requestBody)
			require.NoError(t, err)

			// Send the request and record the response
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			var actualResponseBody types.NetworkFaultInjectionResponse
			err = json.Unmarshal(recorder.Body.Bytes(), &actualResponseBody)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Equal(t, tc.expectedResponseBody, actualResponseBody)

		})
	}
}

func TestCheckNetworkPacketLoss(t *testing.T) {
	tcs := generateNetworkPacketLossTestCases("check network packet loss", "running")

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// Mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			agentState := mock_state.NewMockAgentState(ctrl)
			metricsFactory := mock_metrics.NewMockEntryFactory(ctrl)
			execWrapper := mock_execwrapper.NewMockExec(ctrl)

			if tc.setAgentStateExpectations != nil {
				tc.setAgentStateExpectations(agentState)
			}

			router := mux.NewRouter()
			handler := New(agentState, metricsFactory, execWrapper)
			router.HandleFunc(
				NetworkFaultPath(types.PacketLossFaultType),
				handler.CheckNetworkPacketLoss(),
			).Methods("GET")

			method := "GET"
			var requestBody io.Reader
			if tc.requestBody != nil {
				reqBodyBytes, err := json.Marshal(tc.requestBody)
				require.NoError(t, err)
				requestBody = bytes.NewReader(reqBodyBytes)
			}
			req, err := http.NewRequest(method, fmt.Sprintf("/api/%s/fault/v1/network-packet-loss", endpointId),
				requestBody)
			require.NoError(t, err)

			// Send the request and record the response
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			var actualResponseBody types.NetworkFaultInjectionResponse
			err = json.Unmarshal(recorder.Body.Bytes(), &actualResponseBody)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Equal(t, tc.expectedResponseBody, actualResponseBody)

		})
	}
}
