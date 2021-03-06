// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//

// Code generated by MockGen. DO NOT EDIT.
// Source: k8s.io/client-go/kubernetes/typed/core/v1 (interfaces: NamespaceInterface)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	v1 "k8s.io/api/core/v1"
	v10 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	reflect "reflect"
)

// MockNamespaceInterface is a mock of NamespaceInterface interface
type MockNamespaceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockNamespaceInterfaceMockRecorder
}

// MockNamespaceInterfaceMockRecorder is the mock recorder for MockNamespaceInterface
type MockNamespaceInterfaceMockRecorder struct {
	mock *MockNamespaceInterface
}

// NewMockNamespaceInterface creates a new mock instance
func NewMockNamespaceInterface(ctrl *gomock.Controller) *MockNamespaceInterface {
	mock := &MockNamespaceInterface{ctrl: ctrl}
	mock.recorder = &MockNamespaceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNamespaceInterface) EXPECT() *MockNamespaceInterfaceMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockNamespaceInterface) Create(arg0 *v1.Namespace) (*v1.Namespace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(*v1.Namespace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockNamespaceInterfaceMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockNamespaceInterface)(nil).Create), arg0)
}

// Delete mocks base method
func (m *MockNamespaceInterface) Delete(arg0 string, arg1 *v10.DeleteOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockNamespaceInterfaceMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockNamespaceInterface)(nil).Delete), arg0, arg1)
}

// Finalize mocks base method
func (m *MockNamespaceInterface) Finalize(arg0 *v1.Namespace) (*v1.Namespace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Finalize", arg0)
	ret0, _ := ret[0].(*v1.Namespace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Finalize indicates an expected call of Finalize
func (mr *MockNamespaceInterfaceMockRecorder) Finalize(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Finalize", reflect.TypeOf((*MockNamespaceInterface)(nil).Finalize), arg0)
}

// Get mocks base method
func (m *MockNamespaceInterface) Get(arg0 string, arg1 v10.GetOptions) (*v1.Namespace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*v1.Namespace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockNamespaceInterfaceMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockNamespaceInterface)(nil).Get), arg0, arg1)
}

// List mocks base method
func (m *MockNamespaceInterface) List(arg0 v10.ListOptions) (*v1.NamespaceList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].(*v1.NamespaceList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockNamespaceInterfaceMockRecorder) List(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockNamespaceInterface)(nil).List), arg0)
}

// Patch mocks base method
func (m *MockNamespaceInterface) Patch(arg0 string, arg1 types.PatchType, arg2 []byte, arg3 ...string) (*v1.Namespace, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Patch", varargs...)
	ret0, _ := ret[0].(*v1.Namespace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Patch indicates an expected call of Patch
func (mr *MockNamespaceInterfaceMockRecorder) Patch(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Patch", reflect.TypeOf((*MockNamespaceInterface)(nil).Patch), varargs...)
}

// Update mocks base method
func (m *MockNamespaceInterface) Update(arg0 *v1.Namespace) (*v1.Namespace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(*v1.Namespace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockNamespaceInterfaceMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockNamespaceInterface)(nil).Update), arg0)
}

// UpdateStatus mocks base method
func (m *MockNamespaceInterface) UpdateStatus(arg0 *v1.Namespace) (*v1.Namespace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", arg0)
	ret0, _ := ret[0].(*v1.Namespace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateStatus indicates an expected call of UpdateStatus
func (mr *MockNamespaceInterfaceMockRecorder) UpdateStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockNamespaceInterface)(nil).UpdateStatus), arg0)
}

// Watch mocks base method
func (m *MockNamespaceInterface) Watch(arg0 v10.ListOptions) (watch.Interface, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Watch", arg0)
	ret0, _ := ret[0].(watch.Interface)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Watch indicates an expected call of Watch
func (mr *MockNamespaceInterfaceMockRecorder) Watch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockNamespaceInterface)(nil).Watch), arg0)
}
