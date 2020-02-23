// Copyright (C) 2020 The kubeauth Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctl_test

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/codeactual/kubeauth/internal/cage/cli/handler"
	mock_core "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/core/mock"
	cage_k8s_identity "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/identity"
	mock_identity "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/identity/mock"
	cage_exec "github.com/codeactual/kubeauth/internal/cage/os/exec"
	cage_gomock "github.com/codeactual/kubeauth/internal/cage/testkit/gomock"
	"github.com/codeactual/kubeauth/internal/testkit"
)

// HandlerKit provides command test cases with data and mock-setup boilerplate.
//
// It integrates thc HandlerKit type from the internal/testkit package for additional
// command-agnostic boilerplate.
type HandlerKit struct {
	*testkit.HandlerKit

	// ConfigParsed is true if the Parse call should be mocked during Finish.
	ConfigParsed bool

	// NamespaceValidated is true if the Get call should be mocked during Finish.
	NamespaceValidated bool
}

func NewHandlerKit(t *testing.T) *HandlerKit {
	return &HandlerKit{
		HandlerKit:         testkit.NewHandlerKit(t),
		ConfigParsed:       true,
		NamespaceValidated: true,
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

	if k.ConfigParsed {
		k.ConfigClient.EXPECT().
			Parse("").
			Return(testkit.NewConfigFile(testkit.ConfigFilename, context, cluster, namespace), nil)
	}

	if k.NamespaceValidated {
		k.ApiClientset.Namespaces.EXPECT().
			Get(namespace).
			Return(testkit.NewNamespace(namespace), testkit.Exists, nil)
	}

	// Prepare the mock CLI session, such as to expect specific error message content.

	if k.ExitOnErr != nil {
		k.Session.EXPECT().ExitOnErr(cage_gomock.ErrShortRegexp(k.ExitOnErr), "", 1)
	}
}

// UserQuery configures the kit to expect an --as query with the input user and results.
func (k *HandlerKit) UserQuery(allNamspaces bool, context, cluster, namespace, username string, resultset testkit.QueryResultset) {
	configFile := testkit.NewConfigFile(testkit.ConfigFilename, context, cluster, namespace)

	var queryNamespace string
	if !allNamspaces {
		queryNamespace = namespace
	}

	query := mock_identity.MatchQuery(allNamspaces, &cage_k8s_identity.Query{
		Kind: "User", Namespace: queryNamespace, Name: username, ClientCmdConfig: &configFile.ClientCmdConfig,
	})

	expectClientset := k.ApiClientset.ToReal()
	k.IdentityRegistry.CoreUser.EXPECT().
		Do(cage_gomock.ContextNonNil(), mock_core.MatchClientset(expectClientset), query).
		Return(resultset.CoreUser, nil)
	k.IdentityRegistry.RoleSubject.EXPECT().
		Do(cage_gomock.ContextNonNil(), mock_core.MatchClientset(expectClientset), query).
		Return(resultset.RoleSubject, nil)
	k.IdentityRegistry.ClusterRoleSubject.EXPECT().
		Do(cage_gomock.ContextNonNil(), mock_core.MatchClientset(expectClientset), query).
		Return(resultset.ClusterRoleSubject, nil)
	k.IdentityRegistry.ServiceAccountUser.EXPECT().
		Do(cage_gomock.ContextNonNil(), mock_core.MatchClientset(expectClientset), query).
		Return(resultset.ServiceAccountUser, nil)
	k.IdentityRegistry.ConfigUser.EXPECT().
		Do(cage_gomock.ContextNonNil(), mock_core.MatchClientset(expectClientset), query).
		Return(resultset.ConfigUser, nil)
}

// UserQueryWithDefaultFlags is a convenience alternative to UserQuery when no explicit
// flags for context/cluster/namespace are used and the current ones are expected to be applied
// by the handler as defaults.
func (k *HandlerKit) UserQueryWithDefaultFlags(username string, resultset testkit.QueryResultset) {
	k.UserQuery(testkit.AllNamspacesEnabled, testkit.CurrentContextName, testkit.CurrentClusterName, testkit.CurrentNamespace, username, resultset)
}

// GroupQuery configures the kit to expect an --as-group query with the input group and results.
func (k *HandlerKit) GroupQuery(allNamspaces bool, context, cluster, namespace, group string, resultset testkit.QueryResultset) {
	configFile := testkit.NewConfigFile(testkit.ConfigFilename, context, cluster, namespace)

	var queryNamespace string
	if !allNamspaces {
		queryNamespace = namespace
	}

	query := mock_identity.MatchQuery(allNamspaces, &cage_k8s_identity.Query{
		Kind: "Group", Namespace: queryNamespace, Name: group, ClientCmdConfig: &configFile.ClientCmdConfig,
	})

	expectClientset := k.ApiClientset.ToReal()
	k.IdentityRegistry.CoreGroup.EXPECT().
		Do(cage_gomock.ContextNonNil(), mock_core.MatchClientset(expectClientset), query).
		Return(resultset.CoreGroup, nil)
	k.IdentityRegistry.RoleSubject.EXPECT().
		Do(cage_gomock.ContextNonNil(), mock_core.MatchClientset(expectClientset), query).
		Return(resultset.RoleSubject, nil)
	k.IdentityRegistry.ClusterRoleSubject.EXPECT().
		Do(cage_gomock.ContextNonNil(), mock_core.MatchClientset(expectClientset), query).
		Return(resultset.ClusterRoleSubject, nil)
	k.IdentityRegistry.ServiceAccountGroup.EXPECT().
		Do(cage_gomock.ContextNonNil(), mock_core.MatchClientset(expectClientset), query).
		Return(resultset.ServiceAccountGroup, nil)
}

// GroupQueryWithDefaultFlags is a convenience alternative to GroupQuery when no explicit
// flags for context/cluster/namespace are used and the current ones are expected to be applied
// by the handler as defaults.
func (k *HandlerKit) GroupQueryWithDefaultFlags(group string, resultset testkit.QueryResultset) {
	k.GroupQuery(testkit.AllNamspacesEnabled, testkit.CurrentContextName, testkit.CurrentClusterName, testkit.CurrentNamespace, group, resultset)
}

// StandardCommand configures the kit to expect a command to be created and executed
// with the cage_exec.Executor.Standard method.
//
// Use interface{} as the type of "args" to avoid the "cannot use args[1:] (type []string) as type
// []interface {} in argument" error from gomock (v1.3.1) if the string type is used in the expectations for
// the cage_exec.Executor.Command call..
func (k *HandlerKit) StandardCommand(name string, args ...interface{}) {
	args = append(args, "--v", "1")

	expectCmd := &exec.Cmd{}

	k.Executor.EXPECT().
		Command(name, args...).
		Return(expectCmd)

	expectCmdRes := cage_exec.PipelineResult{}

	k.Executor.EXPECT().
		Standard(
			testkit.Ctx(),
			k.Stdout,
			k.Stderr,
			nil,
			expectCmd,
		).Return(expectCmdRes, nil)
}

func RequireUserQuery(t *testing.T, allNamspaces bool, context, cluster, namespace, username string, resultset testkit.QueryResultset) (stdout, stderr *bytes.Buffer) {
	// Configure the mock objects.

	stdout, stderr = &bytes.Buffer{}, &bytes.Buffer{}

	kit := NewHandlerKit(t)
	kit.Stdout = stdout
	kit.Stderr = stderr

	kit.UserQuery(allNamspaces, context, cluster, namespace, username, resultset)

	// Expected kubectl invocation.

	var args []interface{}
	args = append(
		args,
		"auth", "can-i",
		"--kubeconfig", testkit.ConfigFilename,
		"--as", username,
	)

	if cluster != testkit.CurrentClusterName {
		args = append(args, "--cluster", cluster)
		kit.ClusterName = cluster
	}
	if context != testkit.CurrentContextName {
		args = append(args, "--context", context)
		kit.ContextName = context
	}
	if allNamspaces {
		args = append(args, "--all-namespaces")
		kit.NamespaceValidated = false
	} else {
		if namespace != testkit.CurrentNamespace {
			args = append(args, "--namespace", namespace)
			kit.Namespace = namespace
		}
	}

	kit.StandardCommand("kubectl", args...)
	kit.Finish()
	defer kit.MockCtrl.Finish()

	// Run the CLI handler.

	h := NewHandler(kit)
	h.As = username
	if cluster != testkit.CurrentClusterName {
		h.Cluster = cluster
	}
	if context != testkit.CurrentContextName {
		h.Context = context
	}
	if allNamspaces {
		h.AllNamespaces = true
	} else {
		if namespace != testkit.CurrentNamespace {
			h.Namespace = namespace
		}
	}
	h.Run(testkit.Ctx(), handler.Input{ArgsBeforeDash: []string{"auth", "can-i"}})

	return stdout, stderr
}

func RequireUserQueryWithDefaultFlags(t *testing.T, username string, resultset testkit.QueryResultset) (stdout, stderr *bytes.Buffer) {
	// Configure the mock objects.

	stdout, stderr = &bytes.Buffer{}, &bytes.Buffer{}

	kit := NewHandlerKit(t)
	kit.Stdout = stdout
	kit.Stderr = stderr

	kit.UserQueryWithDefaultFlags(username, resultset)

	// Expected kubectl invocation.
	kit.StandardCommand(
		"kubectl", "auth", "can-i",
		"--kubeconfig", testkit.ConfigFilename,
		"--as", username,
	)
	kit.Finish()
	defer kit.MockCtrl.Finish()

	// Run the CLI handler.

	h := NewHandler(kit)
	h.As = username
	h.Run(testkit.Ctx(), handler.Input{ArgsBeforeDash: []string{"auth", "can-i"}})

	return stdout, stderr
}

func RequireGroupQueryWithDefaultFlags(t *testing.T, group string, resultset testkit.QueryResultset) (stdout, stderr *bytes.Buffer) {
	// Configure the mock objects.

	stdout, stderr = &bytes.Buffer{}, &bytes.Buffer{}

	kit := NewHandlerKit(t)
	kit.Stdout = stdout
	kit.Stderr = stderr

	kit.GroupQueryWithDefaultFlags(group, resultset)

	// Expected kubectl invocation.
	kit.StandardCommand(
		"kubectl", "auth", "can-i",
		"--kubeconfig", testkit.ConfigFilename,
		"--as-group", group,
	)
	kit.Finish()
	defer kit.MockCtrl.Finish()

	// Run the CLI handler.

	h := NewHandler(kit)
	h.AsGroup = []string{group}
	h.Run(testkit.Ctx(), handler.Input{ArgsBeforeDash: []string{"auth", "can-i"}})

	return stdout, stderr
}
