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
	gomock "github.com/golang/mock/gomock"
	v1 "k8s.io/api/core/v1"
	v10 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// CreateBasic mocks base method
func (m *MockClient) CreateBasic(ns, sa string) (*v1.ServiceAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBasic", ns, sa)
	ret0, _ := ret[0].(*v1.ServiceAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBasic indicates an expected call of CreateBasic
func (mr *MockClientMockRecorder) CreateBasic(ns, sa interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBasic", reflect.TypeOf((*MockClient)(nil).CreateBasic), ns, sa)
}

// Get mocks base method
func (m *MockClient) Get(ns, sa string, options ...v10.GetOptions) (*v1.ServiceAccount, bool, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ns, sa}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*v1.ServiceAccount)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get
func (mr *MockClientMockRecorder) Get(ns, sa interface{}, options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ns, sa}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockClient)(nil).Get), varargs...)
}

// List mocks base method
func (m *MockClient) List(ns string, options ...v10.ListOptions) (*v1.ServiceAccountList, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ns}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].(*v1.ServiceAccountList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockClientMockRecorder) List(ns interface{}, options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ns}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockClient)(nil).List), varargs...)
}
