// Copyright (C) 2020 The kubeauth Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package mock

import (
	"github.com/golang/mock/gomock"

	cage_k8s_core "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/core"
	mock_namespace "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/namespace/mock"
	mock_cluster_role "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/cluster_role/mock"
	mock_cluster_role_binding "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/cluster_role_binding/mock"
	mock_role "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/role/mock"
	mock_role_binding "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/role_binding/mock"
	mock_secret "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/secret/mock"
	mock_service_account "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/service_account/mock"
)

// Clientset fields mirror the non-mock Clientset so the latter's values can be replaced.
type Clientset struct {
	ClusterRoles        *mock_cluster_role.MockClient
	ClusterRoleBindings *mock_cluster_role_binding.MockClient
	Namespaces          *mock_namespace.MockClient
	Roles               *mock_role.MockClient
	RoleBindings        *mock_role_binding.MockClient
	Secrets             *mock_secret.MockClient
	ServiceAccounts     *mock_service_account.MockClient
}

func (c *Clientset) ToReal() *cage_k8s_core.Clientset {
	return &cage_k8s_core.Clientset{
		ClusterRoles:        c.ClusterRoles,
		ClusterRoleBindings: c.ClusterRoleBindings,
		Namespaces:          c.Namespaces,
		Roles:               c.Roles,
		RoleBindings:        c.RoleBindings,
		Secrets:             c.Secrets,
		ServiceAccounts:     c.ServiceAccounts,
	}
}

func NewClientset(ctrl *gomock.Controller) *Clientset {
	return &Clientset{
		ClusterRoles:        mock_cluster_role.NewMockClient(ctrl),
		ClusterRoleBindings: mock_cluster_role_binding.NewMockClient(ctrl),
		Namespaces:          mock_namespace.NewMockClient(ctrl),
		Roles:               mock_role.NewMockClient(ctrl),
		RoleBindings:        mock_role_binding.NewMockClient(ctrl),
		Secrets:             mock_secret.NewMockClient(ctrl),
		ServiceAccounts:     mock_service_account.NewMockClient(ctrl),
	}
}

type matchClientset struct {
	expected *cage_k8s_core.Clientset
}

func (m *matchClientset) Matches(x interface{}) bool {
	actual, ok := x.(*cage_k8s_core.Clientset)

	return ok && actual != nil &&
		gomock.Eq(m.expected.ClusterRoles).Matches(actual.ClusterRoles) &&
		gomock.Eq(m.expected.ClusterRoleBindings).Matches(actual.ClusterRoleBindings) &&
		gomock.Eq(m.expected.Namespaces).Matches(actual.Namespaces) &&
		gomock.Eq(m.expected.Roles).Matches(actual.Roles) &&
		gomock.Eq(m.expected.RoleBindings).Matches(actual.RoleBindings) &&
		gomock.Eq(m.expected.Secrets).Matches(actual.Secrets) &&
		gomock.Eq(m.expected.ServiceAccounts).Matches(actual.ServiceAccounts)
}

func (m *matchClientset) String() string {
	return "Clientset.Items match"
}

var _ gomock.Matcher = (*matchClientset)(nil)

// MatchClientset returns a matcher to assert the equality of the individual clients.
func MatchClientset(expected *cage_k8s_core.Clientset) gomock.Matcher { // must return this type for gomock to recognize it
	return &matchClientset{expected: expected}
}
