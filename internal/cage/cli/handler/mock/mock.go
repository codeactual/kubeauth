// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package mock

import (
	io "io"
	"os"
	"regexp"

	"github.com/golang/mock/gomock"

	mock_exec "github.com/codeactual/kubeauth/internal/cage/os/exec/mock"
)

// HandlerKit provides an embeddable starting point for sub-commands' own HandlerKit
// implementations to build on, providing setup for common mocks via NewHandlerKit.
//
// Embedders must call Finish.
type HandlerKit struct {
	Executor *mock_exec.MockExecutor
	MockCtrl *gomock.Controller
	Session  *MockSession

	// ExitOnErr defines the expected pattern for an ExitOnErr/ExitOnErrShort call.
	//
	// After NewHandlerKit processes all HandlerKitOption functions, it configures the related
	// mock if the pattern is non-nil.
	ExitOnErr *regexp.Regexp

	// stdout if non-nil will be used by the Handler for session instead of os.Stdout.
	Stdout io.Writer

	// stderr if non-nil will be used by the Handler for session instead of os.Stderr.
	Stderr io.Writer
}

func NewHandlerKit(mockCtrl *gomock.Controller) *HandlerKit {
	return &HandlerKit{
		Executor: mock_exec.NewMockExecutor(mockCtrl),
		MockCtrl: mockCtrl,
		Session:  NewMockSession(mockCtrl),
		Stderr:   os.Stderr,
		Stdout:   os.Stdout,
	}
}

// Finish creates the expected calls, based on mock-related HandlerKit fields, that were not
// already created by other methods.
func (k *HandlerKit) Finish() {
	k.Session.EXPECT().Err().Return(k.Stderr).AnyTimes()
	k.Session.EXPECT().Out().Return(k.Stdout).AnyTimes()
}
