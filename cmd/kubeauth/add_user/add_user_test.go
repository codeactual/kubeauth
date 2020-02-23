// Copyright (C) 2020 The kubeauth Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package add_user_test asserts CLI behavior by running the command handler logic
// directly (w/o separate processes) with various input scenarios.
//
// It uses Handler instances that use mock implementations of the clients used
// to modify kubeconfig files and perform API requests. The tests only verify correct
// use of the client interfaces. Tests in the cage_k8s package tree verify
// lower-level client behaviors.
//
// It defines the test cases in add_user_test.go. The test cases then rely on
// HandlerKit in handler_kit_test.go to provide common mock boilerplate.
//
// It relies on the internal/testkit package for test fixture values and other
// command-agnotic boilerplate.
package add_user_test

import (
	"regexp"
	"testing"

	"github.com/golang/mock/gomock"
	rbac "k8s.io/api/rbac/v1"

	cli "github.com/codeactual/kubeauth/cmd/kubeauth/add_user"
	"github.com/codeactual/kubeauth/internal/cage/cli/handler"
	cage_k8s "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1"
	cage_gomock "github.com/codeactual/kubeauth/internal/cage/testkit/gomock"
	"github.com/codeactual/kubeauth/internal/testkit"
)

func NewHandler(kit *HandlerKit) *cli.Handler {
	h := cli.Handler{
		Session:             kit.Session,
		KubectlConfigClient: kit.ConfigClient,
		KubeApiClientset:    kit.ApiClientset.ToReal(),
	}

	// Set required CLI flags whose specific values are not yet a SUT.
	h.Username = testkit.Username
	h.ServiceAccountName = kit.ServiceAccountName

	// Enable for test troubleshooting and verbose output assertions.
	h.Verbosity = 1

	return &h
}

// TestApplyDefaultFlags asserts that if these flags are missing, default values are applied:
//   --cluster, --kubeconfig, --namespace
//
// See NewHandlerKit for how the default values of the flags are applied if the respective
// Handler fields are empty when Run executes.
func TestApplyDefaultFlags(t *testing.T) {
	kit := NewHandlerKit(t)
	kit.Namespace = testkit.CurrentNamespace
	kit.ExpectCreatedServiceAccount(kit.Namespace, testkit.ServiceAccountName)
	kit.Finish()
	defer kit.MockCtrl.Finish()

	h := NewHandler(kit)
	h.Run(testkit.Ctx(), handler.Input{})
}

// TestExistingServiceAccount asserts the user may select an --account which already exists.
func TestExistingServiceAccount(t *testing.T) {
	kit := NewHandlerKit(t)
	kit.Namespace = testkit.CurrentNamespace
	kit.ExpectExistingServiceAccount()
	kit.Finish()
	defer kit.MockCtrl.Finish()

	h := NewHandler(kit)
	h.ServiceAccountName = testkit.ServiceAccountName
	h.Run(testkit.Ctx(), handler.Input{})
}

// TestApplyExplicitCluster asserts that an explicit --cluster selection is applied.
func TestApplyExplicitCluster(t *testing.T) {
	explicit := "some-cluster"

	kit := NewHandlerKit(t)
	kit.Namespace = testkit.CurrentNamespace
	kit.ExpectCreatedServiceAccount(kit.Namespace, testkit.ServiceAccountName)
	kit.ClusterName = explicit
	kit.Finish()
	defer kit.MockCtrl.Finish()

	h := NewHandler(kit)
	h.Cluster = explicit
	h.Run(testkit.Ctx(), handler.Input{})
}

// TestApplyExplicitNamespace asserts that an explicit --namespace selection is applied.
func TestApplyExplicitNamespace(t *testing.T) {
	explicit := "some-namespace"

	kit := NewHandlerKit(t)
	kit.Namespace = explicit
	kit.ExpectCreatedServiceAccount(explicit, testkit.ServiceAccountName)
	kit.Finish()
	defer kit.MockCtrl.Finish()

	h := NewHandler(kit)
	h.Namespace = explicit
	h.Run(testkit.Ctx(), handler.Input{})
}

// TestErrOnRoleNotFound asserts that the CLI exists with an error if a role in a
// "--role <role name>:<binding name>" selection does not exist.
func TestErrOnRoleNotFound(t *testing.T) {
	kit := NewHandlerKit(t)
	kit.ExitOnErr = regexp.MustCompile(`role\(s\) not found:.*invalid-a.*invalid-b`)
	kit.UpsertToken = false
	kit.UpsertContext = false
	kit.SecretGet = false
	kit.Finish()
	defer kit.MockCtrl.Finish()

	h := NewHandler(kit)
	h.Roles = []string{"invalid-a:bind-a", "invalid-b:bind-b"}

	kit.ApiClientset.Roles.EXPECT().
		Get(gomock.Any(), "invalid-a").
		Return(nil, testkit.NotExists, nil)
	kit.ApiClientset.Roles.EXPECT().
		Get(gomock.Any(), "invalid-b").
		Return(nil, testkit.NotExists, nil)

	h.Run(testkit.Ctx(), handler.Input{})
}

// TestErrOnClusterRoleNotFound asserts that the CLI exists with an error if a role in a
// "--cluster-role <role name>:<binding name>" selection does not exist.
func TestErrOnClusterRoleNotFound(t *testing.T) {
	kit := NewHandlerKit(t)
	kit.ExitOnErr = regexp.MustCompile(`cluster role\(s\) not found:.*invalid-a.*invalid-b`)
	kit.UpsertToken = false
	kit.UpsertContext = false
	kit.SecretGet = false
	kit.Finish()
	defer kit.MockCtrl.Finish()

	h := NewHandler(kit)
	h.ClusterRoles = []string{"invalid-a:bind-a", "invalid-b:bind-b"}

	kit.ApiClientset.ClusterRoles.EXPECT().
		Get("invalid-a").
		Return(nil, testkit.NotExists, nil)
	kit.ApiClientset.ClusterRoles.EXPECT().
		Get("invalid-b").
		Return(nil, testkit.NotExists, nil)

	h.Run(testkit.Ctx(), handler.Input{})
}

// TestCreateRoleBinding asserts that a --role selection creates the binding.
func TestCreateRoleBinding(t *testing.T) {
	kit := NewHandlerKit(t)
	kit.Namespace = testkit.CurrentNamespace
	kit.ExpectCreatedServiceAccount(kit.Namespace, testkit.ServiceAccountName)
	kit.Finish()
	defer kit.MockCtrl.Finish()

	roleNames := []string{"role-a", "role-b"}
	bindNames := []string{"bind-a", "bind-b"}
	subjects := []rbac.Subject{
		{Kind: cage_k8s.KindServiceAccount, Name: kit.ServiceAccountName},
		{Kind: cage_k8s.KindServiceAccount, Name: kit.ServiceAccountName},
	}

	for n := 0; n < 2; n++ {
		// expect: role name validated
		kit.ApiClientset.Roles.EXPECT().
			Get(kit.Namespace, roleNames[n]).
			Return(cage_gomock.NonSut(), testkit.Exists, nil)

		// expect: binding created
		kit.ApiClientset.RoleBindings.EXPECT().
			Create(kit.Namespace, bindNames[n], roleNames[n], subjects[n]).
			Return(cage_gomock.NonSut(), nil)
	}

	h := NewHandler(kit)
	h.Roles = []string{roleNames[0] + ":" + bindNames[0], roleNames[1] + ":" + bindNames[1]}
	h.Run(testkit.Ctx(), handler.Input{})
}

// TestCreateClusterRoleBinding asserts that a --cluster-role selection creates the binding.
func TestCreateClusterRoleBinding(t *testing.T) {
	kit := NewHandlerKit(t)
	kit.Namespace = testkit.CurrentNamespace
	kit.ExpectCreatedServiceAccount(kit.Namespace, testkit.ServiceAccountName)
	kit.Finish()
	defer kit.MockCtrl.Finish()

	roleNames := []string{"role-a", "role-b"}
	bindNames := []string{"bind-a", "bind-b"}
	subjects := []rbac.Subject{
		{Namespace: kit.Namespace, Kind: cage_k8s.KindServiceAccount, Name: kit.ServiceAccountName},
		{Namespace: kit.Namespace, Kind: cage_k8s.KindServiceAccount, Name: kit.ServiceAccountName},
	}

	for n := 0; n < 2; n++ {
		// expect: role name validated
		kit.ApiClientset.ClusterRoles.EXPECT().
			Get(roleNames[n]).
			Return(cage_gomock.NonSut(), testkit.Exists, nil)

		// expect: binding created
		kit.ApiClientset.ClusterRoleBindings.EXPECT().
			Create(bindNames[n], roleNames[n], subjects[n]).
			Return(cage_gomock.NonSut(), nil)
	}

	h := NewHandler(kit)
	h.ClusterRoles = []string{roleNames[0] + ":" + bindNames[0], roleNames[1] + ":" + bindNames[1]}
	h.Run(testkit.Ctx(), handler.Input{})
}
