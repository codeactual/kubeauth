// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

//go:generate mockgen -copyright_file=$LICENSE_HEADER -package=mock -destination=$GODIR/mock/cluster_role_binding.go -source=$GODIR/$GOFILE
//go:generate mockgen -copyright_file $CAPATH/LICENSE_HEADER -package=mock -destination=$GODIR/mock/getter.go k8s.io/client-go/kubernetes/typed/rbac/v1 ClusterRoleBindingsGetter
//go:generate mockgen -copyright_file $CAPATH/LICENSE_HEADER -package=mock -destination=$GODIR/mock/interface.go k8s.io/client-go/kubernetes/typed/rbac/v1 ClusterRoleBindingInterface
package cluster_role_binding

import (
	rbac "k8s.io/api/rbac/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbac_type "k8s.io/client-go/kubernetes/typed/rbac/v1"

	"github.com/pkg/errors"

	cage_k8s "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1"
)

// Client provides an interface to cluster role bindings.
type Client interface {
	Create(name, role string, subject rbac.Subject) (*rbac.ClusterRoleBinding, error)
	List(options ...meta.ListOptions) (*rbac.ClusterRoleBindingList, error)
}

// DefaultClient implementation of Client operates on a real kubernetes API.
type DefaultClient struct {
	rbac_type.ClusterRoleBindingsGetter
}

// NewDefaultClient returns an initialized DefaultClient.
func NewDefaultClient(getter rbac_type.ClusterRoleBindingsGetter) *DefaultClient {
	return &DefaultClient{ClusterRoleBindingsGetter: getter}
}

// Create binds the role to a single subject.
func (c *DefaultClient) Create(name, role string, subject rbac.Subject) (*rbac.ClusterRoleBinding, error) {
	obj, err := c.ClusterRoleBindings().Create(&rbac.ClusterRoleBinding{
		ObjectMeta: meta.ObjectMeta{Name: name},
		RoleRef:    rbac.RoleRef{Kind: cage_k8s.KindClusterRole, Name: role},
		Subjects:   []rbac.Subject{subject},
	})
	if err != nil {
		// Allow caller to perform the same check and decide whether how to handlei it.
		if k8s_errors.IsAlreadyExists(err) {
			return nil, err
		}

		nsId := subject.Namespace
		if nsId == "" {
			nsId = cage_k8s.EmptyNamespace
		}

		return nil, errors.Wrapf(
			err,
			"failed to bind cluster role [%s] to subject [%s] (kind: %s ns: %s)",
			role, subject.Name, subject.Kind, nsId,
		)
	}

	return obj, nil
}

// List returns the matching objects.
//
// A single ListOptions value can be passed as the final argument to customize the query.
//
// It implements Client.
func (c *DefaultClient) List(options ...meta.ListOptions) (*rbac.ClusterRoleBindingList, error) {
	list, err := c.ClusterRoleBindings().List(cage_k8s.ListOptionsFromVariadic(options))
	if err != nil {
		if k8s_errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to list cluster role bindings")
	}

	return list, nil
}

var _ Client = (*DefaultClient)(nil)
