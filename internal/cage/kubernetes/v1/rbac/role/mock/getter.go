// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//

// Code generated by MockGen. DO NOT EDIT.
// Source: k8s.io/client-go/kubernetes/typed/rbac/v1 (interfaces: RolesGetter)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	v1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	reflect "reflect"
)

// MockRolesGetter is a mock of RolesGetter interface
type MockRolesGetter struct {
	ctrl     *gomock.Controller
	recorder *MockRolesGetterMockRecorder
}

// MockRolesGetterMockRecorder is the mock recorder for MockRolesGetter
type MockRolesGetterMockRecorder struct {
	mock *MockRolesGetter
}

// NewMockRolesGetter creates a new mock instance
func NewMockRolesGetter(ctrl *gomock.Controller) *MockRolesGetter {
	mock := &MockRolesGetter{ctrl: ctrl}
	mock.recorder = &MockRolesGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRolesGetter) EXPECT() *MockRolesGetterMockRecorder {
	return m.recorder
}

// Roles mocks base method
func (m *MockRolesGetter) Roles(arg0 string) v1.RoleInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Roles", arg0)
	ret0, _ := ret[0].(v1.RoleInterface)
	return ret0
}

// Roles indicates an expected call of Roles
func (mr *MockRolesGetterMockRecorder) Roles(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Roles", reflect.TypeOf((*MockRolesGetter)(nil).Roles), arg0)
}