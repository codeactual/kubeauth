// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//

// Code generated by MockGen. DO NOT EDIT.
// Source: k8s.io/client-go/kubernetes/typed/core/v1 (interfaces: ServiceAccountsGetter)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	reflect "reflect"
)

// MockServiceAccountsGetter is a mock of ServiceAccountsGetter interface
type MockServiceAccountsGetter struct {
	ctrl     *gomock.Controller
	recorder *MockServiceAccountsGetterMockRecorder
}

// MockServiceAccountsGetterMockRecorder is the mock recorder for MockServiceAccountsGetter
type MockServiceAccountsGetterMockRecorder struct {
	mock *MockServiceAccountsGetter
}

// NewMockServiceAccountsGetter creates a new mock instance
func NewMockServiceAccountsGetter(ctrl *gomock.Controller) *MockServiceAccountsGetter {
	mock := &MockServiceAccountsGetter{ctrl: ctrl}
	mock.recorder = &MockServiceAccountsGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockServiceAccountsGetter) EXPECT() *MockServiceAccountsGetterMockRecorder {
	return m.recorder
}

// ServiceAccounts mocks base method
func (m *MockServiceAccountsGetter) ServiceAccounts(arg0 string) v1.ServiceAccountInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServiceAccounts", arg0)
	ret0, _ := ret[0].(v1.ServiceAccountInterface)
	return ret0
}

// ServiceAccounts indicates an expected call of ServiceAccounts
func (mr *MockServiceAccountsGetterMockRecorder) ServiceAccounts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServiceAccounts", reflect.TypeOf((*MockServiceAccountsGetter)(nil).ServiceAccounts), arg0)
}