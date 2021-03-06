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
	bytes "bytes"
	context "context"
	exec "github.com/codeactual/kubeauth/internal/cage/os/exec"
	gomock "github.com/golang/mock/gomock"
	io "io"
	exec0 "os/exec"
	reflect "reflect"
)

// MockExecutor is a mock of Executor interface
type MockExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockExecutorMockRecorder
}

// MockExecutorMockRecorder is the mock recorder for MockExecutor
type MockExecutorMockRecorder struct {
	mock *MockExecutor
}

// NewMockExecutor creates a new mock instance
func NewMockExecutor(ctrl *gomock.Controller) *MockExecutor {
	mock := &MockExecutor{ctrl: ctrl}
	mock.recorder = &MockExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockExecutor) EXPECT() *MockExecutorMockRecorder {
	return m.recorder
}

// Command mocks base method
func (m *MockExecutor) Command(name string, arg ...string) *exec0.Cmd {
	m.ctrl.T.Helper()
	varargs := []interface{}{name}
	for _, a := range arg {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Command", varargs...)
	ret0, _ := ret[0].(*exec0.Cmd)
	return ret0
}

// Command indicates an expected call of Command
func (mr *MockExecutorMockRecorder) Command(name interface{}, arg ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{name}, arg...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Command", reflect.TypeOf((*MockExecutor)(nil).Command), varargs...)
}

// CommandContext mocks base method
func (m *MockExecutor) CommandContext(ctx context.Context, name string, arg ...string) *exec0.Cmd {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, name}
	for _, a := range arg {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CommandContext", varargs...)
	ret0, _ := ret[0].(*exec0.Cmd)
	return ret0
}

// CommandContext indicates an expected call of CommandContext
func (mr *MockExecutorMockRecorder) CommandContext(ctx, name interface{}, arg ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, name}, arg...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommandContext", reflect.TypeOf((*MockExecutor)(nil).CommandContext), varargs...)
}

// Buffered mocks base method
func (m *MockExecutor) Buffered(ctx context.Context, cmds ...*exec0.Cmd) (*bytes.Buffer, *bytes.Buffer, exec.PipelineResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range cmds {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Buffered", varargs...)
	ret0, _ := ret[0].(*bytes.Buffer)
	ret1, _ := ret[1].(*bytes.Buffer)
	ret2, _ := ret[2].(exec.PipelineResult)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// Buffered indicates an expected call of Buffered
func (mr *MockExecutorMockRecorder) Buffered(ctx interface{}, cmds ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, cmds...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Buffered", reflect.TypeOf((*MockExecutor)(nil).Buffered), varargs...)
}

// Standard mocks base method
func (m *MockExecutor) Standard(ctx context.Context, stdout, stderr io.Writer, stdin io.Reader, cmds ...*exec0.Cmd) (exec.PipelineResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, stdout, stderr, stdin}
	for _, a := range cmds {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Standard", varargs...)
	ret0, _ := ret[0].(exec.PipelineResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Standard indicates an expected call of Standard
func (mr *MockExecutorMockRecorder) Standard(ctx, stdout, stderr, stdin interface{}, cmds ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, stdout, stderr, stdin}, cmds...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Standard", reflect.TypeOf((*MockExecutor)(nil).Standard), varargs...)
}

// Pty mocks base method
func (m *MockExecutor) Pty(cmd *exec0.Cmd) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pty", cmd)
	ret0, _ := ret[0].(error)
	return ret0
}

// Pty indicates an expected call of Pty
func (mr *MockExecutorMockRecorder) Pty(cmd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pty", reflect.TypeOf((*MockExecutor)(nil).Pty), cmd)
}
