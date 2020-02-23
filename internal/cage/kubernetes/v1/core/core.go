// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package core

import (
	"k8s.io/client-go/kubernetes"

	cage_k8s_namespace "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/namespace"
	cage_k8s_cluster_role "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/cluster_role"
	cage_k8s_cluster_role_binding "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/cluster_role_binding"
	cage_k8s_role "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/role"
	cage_k8s_role_binding "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/role_binding"
	cage_k8s_secret "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/secret"
	cage_k8s_sa "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/service_account"
)

// Clientset provides customized clients to kubernetes API objects.
//
// Its naming is modeled after k8s.io/client-go/kubernetes.Clientset.
type Clientset struct {
	ClusterRoles        cage_k8s_cluster_role.Client
	ClusterRoleBindings cage_k8s_cluster_role_binding.Client
	Namespaces          cage_k8s_namespace.Client
	Roles               cage_k8s_role.Client
	RoleBindings        cage_k8s_role_binding.Client
	Secrets             cage_k8s_secret.Client
	ServiceAccounts     cage_k8s_sa.Client
}

func NewClientset(all kubernetes.Interface) *Clientset {
	return &Clientset{
		ClusterRoles:        cage_k8s_cluster_role.NewDefaultClient(all.RbacV1()),
		ClusterRoleBindings: cage_k8s_cluster_role_binding.NewDefaultClient(all.RbacV1()),
		Namespaces:          cage_k8s_namespace.NewDefaultClient(all.CoreV1()),
		Roles:               cage_k8s_role.NewDefaultClient(all.RbacV1()),
		RoleBindings:        cage_k8s_role_binding.NewDefaultClient(all.RbacV1()),
		Secrets:             cage_k8s_secret.NewDefaultClient(all.CoreV1()),
		ServiceAccounts:     cage_k8s_sa.NewDefaultClient(all.CoreV1()),
	}
}
