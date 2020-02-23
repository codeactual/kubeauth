// Copyright (C) 2020 The kubeauth Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package add_user_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	core "k8s.io/api/core/v1"

	cage_gomock "github.com/codeactual/kubeauth/internal/cage/testkit/gomock"
	"github.com/codeactual/kubeauth/internal/testkit"
)

// HandlerKit provides command test cases with data and mock-setup boilerplate.
//
// It integrates thc HandlerKit type from the internal/testkit package for additional
// command-agnostic boilerplate.
type HandlerKit struct {
	*testkit.HandlerKit

	// UpsertContext is true if ConfigureMocks should include the context creation/update call.
	UpsertContext bool

	// UpsertContext is true if ConfigureMocks should include the token creation/update call.
	UpsertToken bool

	// SecretGet is true if ConfigureMocks should include the secret creation call.
	SecretGet bool

	// ServiceAccountName is the expected effective value after flag/default processing is complete.
	ServiceAccountName string
}

func NewHandlerKit(t *testing.T) *HandlerKit {
	return &HandlerKit{
		HandlerKit:    testkit.NewHandlerKit(t),
		SecretGet:     true,
		UpsertContext: true,
		UpsertToken:   true,
	}
}

// Finish creates the expected calls, based on mock-related HandlerKit fields, that were not
// already created by other methods.
func (k *HandlerKit) Finish() {
	k.HandlerKit.Finish()

	// Configure mocks to expect these current values if the kit has not selected custom ones.
	cluster, context, namespace := k.ClusterName, k.ContextName, k.Namespace
	if cluster == "" {
		cluster = testkit.CurrentClusterName
	}
	if context == "" {
		context = testkit.CurrentContextName
	}
	if namespace == "" {
		namespace = testkit.CurrentNamespace
	}

	k.ConfigClient.EXPECT().
		Parse("").
		Return(testkit.NewConfigFile(testkit.ConfigFilename, context, cluster, namespace), nil)

	if k.UpsertToken {
		k.ConfigClient.EXPECT().
			UpsertUserToken(testkit.Ctx(), gomock.Any(), testkit.Username, TokenData()).
			Return(nil)
	}

	if k.UpsertContext {
		k.ConfigClient.EXPECT().
			UpsertContext(testkit.Ctx(), gomock.Any(), testkit.Username, cluster, namespace, testkit.Username).
			Return(nil)
	}

	if k.SecretGet {
		expectSecret := &core.Secret{
			Data: map[string][]byte{
				"ca.crt": CertData(),
				"token":  TokenData(),
			},
		}
		k.ApiClientset.Secrets.EXPECT().
			Get(gomock.Any(), SecretName(k.ServiceAccountName)).
			Return(expectSecret, testkit.Exists, nil)
	}

	if k.ExitOnErr != nil {
		k.Session.EXPECT().ExitOnErr(cage_gomock.ErrShortRegexp(k.ExitOnErr), "", 1)
	}
}

// ExistingServiceAccount immediately configures the kit to expect the service account
// and its secret's name already exist.
func (k *HandlerKit) ExpectExistingServiceAccount() {
	existingObj := &core.ServiceAccount{Secrets: []core.ObjectReference{{Name: SecretName(testkit.ServiceAccountName)}}}

	k.ApiClientset.ServiceAccounts.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(existingObj, testkit.Exists, nil)

	k.ServiceAccountName = testkit.ServiceAccountName
}

// CreatedServiceAccount immediately configures the kit to expect the service account
// and its secret's name must be created.
func (k *HandlerKit) ExpectCreatedServiceAccount(namespace, name string) {
	createdObj := &core.ServiceAccount{Secrets: []core.ObjectReference{{Name: SecretName(name)}}}

	gomock.InOrder(
		k.ApiClientset.ServiceAccounts.EXPECT().
			Get(namespace, name).
			Return(createdObj, testkit.NotExists, nil),
		k.ApiClientset.ServiceAccounts.EXPECT().
			CreateBasic(namespace, name).
			Return(createdObj, nil),
	)

	k.ServiceAccountName = name
}
