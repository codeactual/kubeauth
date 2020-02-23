// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//

// Code generated by MockGen. DO NOT EDIT.

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	config "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/kubectl/config"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Parse mocks base method
func (m *MockClient) Parse(filename string) (*config.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parse", filename)
	ret0, _ := ret[0].(*config.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Parse indicates an expected call of Parse
func (mr *MockClientMockRecorder) Parse(filename interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*MockClient)(nil).Parse), filename)
}

// UpsertUserToken mocks base method
func (m *MockClient) UpsertUserToken(ctx context.Context, parsed *config.File, user string, token []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertUserToken", ctx, parsed, user, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertUserToken indicates an expected call of UpsertUserToken
func (mr *MockClientMockRecorder) UpsertUserToken(ctx, parsed, user, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertUserToken", reflect.TypeOf((*MockClient)(nil).UpsertUserToken), ctx, parsed, user, token)
}

// UpsertContext mocks base method
func (m *MockClient) UpsertContext(ctx context.Context, parsed *config.File, name, cluster, ns, user string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertContext", ctx, parsed, name, cluster, ns, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertContext indicates an expected call of UpsertContext
func (mr *MockClientMockRecorder) UpsertContext(ctx, parsed, name, cluster, ns, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertContext", reflect.TypeOf((*MockClient)(nil).UpsertContext), ctx, parsed, name, cluster, ns, user)
}