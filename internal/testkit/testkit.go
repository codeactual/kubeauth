// Copyright (C) 2020 The kubeauth Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package testkit provides constants, functions, and types which provide
// common necessities for the tests in multple sub-commands.
package testkit

import (
	"testing"

	"github.com/golang/mock/gomock"

	mock_handler "github.com/codeactual/kubeauth/internal/cage/cli/handler/mock"
	mock_core "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/core/mock"
	mock_config "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/kubectl/config/mock"
	mock_identity "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/identity/mock"
)

// HandlerKit provides an embeddable starting point for sub-commands' own HandlerKit
// implementations to build on, providing setup for common mocks via NewHandlerKit.
type HandlerKit struct {
	*mock_handler.HandlerKit

	ApiClientset     *mock_core.Clientset
	ConfigClient     *mock_config.MockClient
	IdentityRegistry *mock_identity.Registry

	// ClusterName is the expected effective value after flag/default processing is complete.
	ClusterName string

	// ContextName is the expected effective value after flag/default processing is complete.
	ContextName string

	// Namespace is the expected effective value after flag/default processing is complete.
	Namespace string
}

// NewHandlerKit returns a kit with mocks required by multiple sub-commands
// in order to reduce their boilerplate.
func NewHandlerKit(t *testing.T) *HandlerKit {
	mockCtrl := gomock.NewController(t)
	return &HandlerKit{
		HandlerKit:       mock_handler.NewHandlerKit(mockCtrl),
		ApiClientset:     mock_core.NewClientset(mockCtrl),
		ConfigClient:     mock_config.NewMockClient(mockCtrl),
		IdentityRegistry: mock_identity.NewRegistry(mockCtrl),
	}
}
