// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//

// Code generated by MockGen. DO NOT EDIT.
// Source: k8s.io/client-go/kubernetes/typed/rbac/v1 (interfaces: RoleBindingsGetter)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	v1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	reflect "reflect"
)

// MockRoleBindingsGetter is a mock of RoleBindingsGetter interface
type MockRoleBindingsGetter struct {
	ctrl     *gomock.Controller
	recorder *MockRoleBindingsGetterMockRecorder
}

// MockRoleBindingsGetterMockRecorder is the mock recorder for MockRoleBindingsGetter
type MockRoleBindingsGetterMockRecorder struct {
	mock *MockRoleBindingsGetter
}

// NewMockRoleBindingsGetter creates a new mock instance
func NewMockRoleBindingsGetter(ctrl *gomock.Controller) *MockRoleBindingsGetter {
	mock := &MockRoleBindingsGetter{ctrl: ctrl}
	mock.recorder = &MockRoleBindingsGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRoleBindingsGetter) EXPECT() *MockRoleBindingsGetterMockRecorder {
	return m.recorder
}

// RoleBindings mocks base method
func (m *MockRoleBindingsGetter) RoleBindings(arg0 string) v1.RoleBindingInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RoleBindings", arg0)
	ret0, _ := ret[0].(v1.RoleBindingInterface)
	return ret0
}

// RoleBindings indicates an expected call of RoleBindings
func (mr *MockRoleBindingsGetterMockRecorder) RoleBindings(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RoleBindings", reflect.TypeOf((*MockRoleBindingsGetter)(nil).RoleBindings), arg0)
}
