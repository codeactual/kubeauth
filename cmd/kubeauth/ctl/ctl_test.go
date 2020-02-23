// Copyright (C) 2020 The kubeauth Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package ctl_test asserts CLI behavior by running the command handler logic
// directly (w/o separate processes) with various input scenarios.
//
// It uses Handler instances that use mock implementations of the clients used
// to modify kubeconfig files and perform API requests. The tests only verify correct
// use of the client interfaces. Tests in the cage_k8s package tree verify
// lower-level client behaviors.
//
// It defines the test cases in ctl_test.go. The test cases then rely on
// HandlerKit in handler_kit_test.go to provide common mock boilerplate.
//
// It relies on the internal/testkit package for test fixture values and other
// command-agnotic boilerplate.
package ctl_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	cli "github.com/codeactual/kubeauth/cmd/kubeauth/ctl"
	"github.com/codeactual/kubeauth/internal/cage/cli/handler"
	cage_k8s "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1"
	"github.com/codeactual/kubeauth/internal/testkit"
)

func NewHandler(kit *HandlerKit) *cli.Handler {
	apiClientset := kit.ApiClientset.ToReal()

	h := cli.Handler{
		Session:             kit.Session,
		KubectlConfigClient: kit.ConfigClient,
		KubeApiClientset:    apiClientset,
		Executor:            kit.Executor,
		IdentityRegistry:    kit.IdentityRegistry.ToReal(apiClientset),
	}

	// Enable for test troubleshooting and verbose output assertions.
	h.Verbosity = 1

	return &h
}

// TestErrOnMissingSubject asserts that --as or --as-group are required inputs.
func TestErrOnMissingSubject(t *testing.T) {
	kit := NewHandlerKit(t)
	kit.ExitOnErr = regexp.MustCompile("missing --as or --as-group")
	kit.NamespaceValidated = false
	kit.Finish()
	defer kit.MockCtrl.Finish()

	h := NewHandler(kit)
	h.Run(context.Background(), handler.Input{ArgsBeforeDash: []string{"auth", "can-i"}})
}

// TestErrOnNamespaceScopeConflict asserts that --namespace and --all-namespaces cannot be combined.
func TestErrOnNamespaceScopeConflict(t *testing.T) {
	kit := NewHandlerKit(t)
	kit.ExitOnErr = regexp.MustCompile("missing --as or --as-group")
	kit.NamespaceValidated = false
	kit.Finish()
	defer kit.MockCtrl.Finish()

	h := NewHandler(kit)
	h.AllNamespaces = true
	h.Namespace = testkit.Namespace
	h.Run(context.Background(), handler.Input{ArgsBeforeDash: []string{"auth", "can-i"}})
}

// TestApplyExplicitNamespace asserts that --namespace is applied.
func TestApplyExplicitNamespace(t *testing.T) {
	// Expected query's parameters and results.

	username := testkit.Username
	namespace := testkit.Namespace
	context := testkit.ContextName

	resultset := testkit.NewQueryResultset()
	resultset.ConfigUser.Add(namespace, cage_k8s.KindUser, username, nil)

	// Run the CLI handler.

	_, stderr := RequireUserQuery(t, testkit.AllNamspacesDisabled, context, testkit.CurrentClusterName, namespace, username, resultset)
	require.Contains(t, stderr.String(), "User "+username+" of namespace "+namespace+" via [kubeconfig context] querier")
}

// TestApplyAllNamespaces asserts that --all-namespaces is applied.
func TestApplyAllNamespaces(t *testing.T) {
	// Expected query's parameters and results.

	username := testkit.Username
	namespace := testkit.CurrentNamespace

	resultset := testkit.NewQueryResultset()
	resultset.ConfigUser.Add(namespace, cage_k8s.KindUser, username, nil)

	// Run the CLI handler.

	_, stderr := RequireUserQuery(t, testkit.AllNamspacesEnabled, testkit.CurrentContextName, testkit.CurrentClusterName, namespace, username, resultset)
	require.Contains(t, stderr.String(), "User "+username+" of namespace "+namespace+" via [kubeconfig context] querier")
}

// TestErrOnInvalidContext asserts that the command stops if the current context cannot
// be found in user input or config file.
func TestErrOnInvalidContext(t *testing.T) {
	contextName := "does-not-exist"

	kit := NewHandlerKit(t)
	kit.NamespaceValidated = false
	kit.ExitOnErr = regexp.MustCompile(`context \[` + contextName + `\] not found`)
	kit.Finish()
	defer kit.MockCtrl.Finish()

	h := NewHandler(kit)
	h.As = testkit.Username
	h.Context = contextName
	h.Run(context.Background(), handler.Input{ArgsBeforeDash: []string{"auth", "can-i"}})
}

// TestApplyExplicitContext asserts that --context is applied.
func TestApplyExplicitContext(t *testing.T) {
	// Expected query's parameters and results.

	username := testkit.Username
	namespace := testkit.CurrentNamespace
	contextName := testkit.ContextName

	resultset := testkit.NewQueryResultset()
	resultset.ConfigUser.Add(namespace, cage_k8s.KindUser, username, nil)

	kit := NewHandlerKit(t)
	kit.ContextName = contextName
	kit.UserQuery(testkit.AllNamspacesDisabled, contextName, testkit.CurrentClusterName, namespace, username, resultset)

	// Expected kubectl invocation.
	kit.StandardCommand(
		"kubectl", "auth", "can-i",
		"--kubeconfig", testkit.ConfigFilename,
		"--as", username,
		"--context", contextName,
	)
	kit.Finish()
	defer kit.MockCtrl.Finish()

	// Run the CLI handler.

	h := NewHandler(kit)
	h.As = username
	h.Context = contextName
	h.Run(context.Background(), handler.Input{ArgsBeforeDash: []string{"auth", "can-i"}})
}

// TestErrOnClusterConflict asserts that the command stops if --cluster doesn't match the current context's.
func TestErrOnClusterConflict(t *testing.T) {
	clusterName := "does-not-exist"

	kit := NewHandlerKit(t)
	kit.NamespaceValidated = false
	kit.ExitOnErr = regexp.MustCompile(`cluster \[` + clusterName + `\] differs from effective context's cluster \[` + testkit.CurrentClusterName + `\]`)
	kit.Finish()
	defer kit.MockCtrl.Finish()

	h := NewHandler(kit)
	h.As = testkit.Username
	h.Cluster = clusterName
	h.Run(context.Background(), handler.Input{ArgsBeforeDash: []string{"auth", "can-i"}})
}

// TestApplyExplicitCluster asserts that --cluster is applied.
func TestApplyExplicitCluster(t *testing.T) {
	// Expected query's parameters and results.

	username := testkit.Username
	namespace := testkit.CurrentNamespace
	clusterName := testkit.ClusterName
	contextName := testkit.ContextName // contains the non-current cluster

	resultset := testkit.NewQueryResultset()
	resultset.ConfigUser.Add(namespace, cage_k8s.KindUser, username, nil)

	kit := NewHandlerKit(t)
	kit.ContextName = contextName
	kit.ClusterName = clusterName
	kit.UserQuery(testkit.AllNamspacesDisabled, contextName, clusterName, namespace, username, resultset)

	// Expected kubectl invocation.
	kit.StandardCommand(
		"kubectl", "auth", "can-i",
		"--kubeconfig", testkit.ConfigFilename,
		"--as", username,
		"--context", contextName,
		"--cluster", clusterName,
	)
	kit.Finish()
	defer kit.MockCtrl.Finish()

	// Run the CLI handler.

	h := NewHandler(kit)
	h.As = username
	h.Cluster = clusterName
	h.Context = contextName
	h.Run(context.Background(), handler.Input{ArgsBeforeDash: []string{"auth", "can-i"}})
}

// TestErrOnInvalidUser asserts that the command stops on an invalid --as selection.
func TestErrOnInvalidUser(t *testing.T) {
	username := "does-not-exist"

	kit := NewHandlerKit(t)
	kit.ExitOnErr = regexp.MustCompile(`--as identity \[` + username + `\] not found`)
	kit.UserQuery(testkit.AllNamspacesDisabled, testkit.CurrentContextName, testkit.CurrentClusterName, testkit.CurrentNamespace, username, testkit.NewQueryResultset())
	kit.Finish()
	defer kit.MockCtrl.Finish()

	h := NewHandler(kit)
	h.As = username
	h.Run(context.Background(), handler.Input{ArgsBeforeDash: []string{"auth", "can-i"}})
}

// TestErrOnInvalidGroup asserts that the command stops on an invalid --as-group selection.
func TestErrOnInvalidGroup(t *testing.T) {
	group := "does-not-exist"

	kit := NewHandlerKit(t)
	kit.ExitOnErr = regexp.MustCompile(`--as-group identity \[` + group + `\] not found`)
	kit.GroupQuery(testkit.AllNamspacesDisabled, testkit.CurrentContextName, testkit.CurrentClusterName, testkit.CurrentNamespace, group, testkit.NewQueryResultset())
	kit.Finish()
	defer kit.MockCtrl.Finish()

	h := NewHandler(kit)
	h.AsGroup = []string{group}
	h.Run(context.Background(), handler.Input{ArgsBeforeDash: []string{"auth", "can-i"}})
}

func TestConfigUser(t *testing.T) {
	// Expected query's parameters and results.

	username := testkit.Username
	namespace := testkit.CurrentNamespace

	resultset := testkit.NewQueryResultset()
	resultset.ConfigUser.Add(namespace, cage_k8s.KindUser, username, nil)

	// Run the CLI handler.

	_, stderr := RequireUserQueryWithDefaultFlags(t, username, resultset)
	require.Contains(t, stderr.String(), "User "+username+" of namespace "+namespace+" via [kubeconfig context] querier")
}

func TestCoreUser(t *testing.T) {
	// Expected query's parameters and results.

	username := testkit.Username
	namespace := "" // core users are namespace agnostic

	resultset := testkit.NewQueryResultset()
	resultset.CoreUser.Add(namespace, cage_k8s.KindUser, username, nil)

	// Run the CLI handler.

	_, stderr := RequireUserQueryWithDefaultFlags(t, username, resultset)
	require.Contains(t, stderr.String(), "User "+username+" via [system-defined user] querier")
}

func TestServiceAccountUser(t *testing.T) {
	// Expected query's parameters and results.

	username := testkit.Username
	namespace := testkit.CurrentNamespace

	resultset := testkit.NewQueryResultset()
	resultset.ServiceAccountUser.Add(namespace, cage_k8s.KindUser, username, nil)

	// Run the CLI handler.

	_, stderr := RequireUserQueryWithDefaultFlags(t, username, resultset)
	require.Contains(t, stderr.String(), "User "+username+" of namespace "+namespace+" via [service account based user] querier")
}

func TestRoleSubjectUser(t *testing.T) {
	// Expected query's parameters and results.

	username := testkit.Username
	namespace := testkit.CurrentNamespace

	resultset := testkit.NewQueryResultset()
	resultset.RoleSubject.Add(namespace, cage_k8s.KindUser, username, nil)

	// Run the CLI handler.

	_, stderr := RequireUserQueryWithDefaultFlags(t, username, resultset)
	require.Contains(t, stderr.String(), "User "+username+" of namespace "+namespace+" via [role binding subject] querier")
}

func TestClusterRoleSubjectUser(t *testing.T) {
	// Expected query's parameters and results.

	username := testkit.Username
	namespace := testkit.CurrentNamespace

	resultset := testkit.NewQueryResultset()
	resultset.ClusterRoleSubject.Add(namespace, cage_k8s.KindUser, username, nil)

	// Run the CLI handler.

	_, stderr := RequireUserQueryWithDefaultFlags(t, username, resultset)
	require.Contains(t, stderr.String(), "User "+username+" of namespace "+namespace+" via [cluster role binding subject] querier")
}

func TestCoreGroup(t *testing.T) {
	// Expected query's parameters and results.

	group := testkit.GroupName
	namespace := "" // core users are namespace agnostic

	resultset := testkit.NewQueryResultset()
	resultset.CoreGroup.Add(namespace, cage_k8s.KindGroup, group, nil)

	// Run the CLI handler.

	_, stderr := RequireGroupQueryWithDefaultFlags(t, group, resultset)
	require.Contains(t, stderr.String(), "Group "+group+" via [system-defined group] querier")
}

func TestServiceAccountGroup(t *testing.T) {
	// Expected query's parameters and results.

	group := testkit.GroupName
	namespace := testkit.CurrentNamespace // only for consistency with other cases, core users are namespace agnostic

	resultset := testkit.NewQueryResultset()
	resultset.ServiceAccountGroup.Add(namespace, cage_k8s.KindGroup, group, nil)

	// Run the CLI handler.

	_, stderr := RequireGroupQueryWithDefaultFlags(t, group, resultset)
	require.Contains(t, stderr.String(), "Group "+group+" of namespace "+namespace+" via [service account based group] querier")
}

func TestRoleSubjectGroup(t *testing.T) {
	// Expected query's parameters and results.

	group := testkit.GroupName
	namespace := testkit.CurrentNamespace

	resultset := testkit.NewQueryResultset()
	resultset.RoleSubject.Add(namespace, cage_k8s.KindGroup, group, nil)

	// Run the CLI handler.

	_, stderr := RequireGroupQueryWithDefaultFlags(t, group, resultset)
	require.Contains(t, stderr.String(), "Group "+group+" of namespace "+namespace+" via [role binding subject] querier")
}

func TestClusterRoleSubjectGroup(t *testing.T) {
	// Expected query's parameters and results.

	group := testkit.GroupName
	namespace := testkit.CurrentNamespace

	resultset := testkit.NewQueryResultset()
	resultset.ClusterRoleSubject.Add(namespace, cage_k8s.KindGroup, group, nil)

	// Run the CLI handler.

	_, stderr := RequireGroupQueryWithDefaultFlags(t, group, resultset)
	require.Contains(t, stderr.String(), "Group "+group+" of namespace "+namespace+" via [cluster role binding subject] querier")
}
