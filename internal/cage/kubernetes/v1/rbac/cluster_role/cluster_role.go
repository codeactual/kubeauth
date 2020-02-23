// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

//go:generate mockgen -copyright_file=$LICENSE_HEADER -package=mock -destination=$GODIR/mock/cluster_role.go -source=$GODIR/$GOFILE
//go:generate mockgen -copyright_file $CAPATH/LICENSE_HEADER -package=mock -destination=$GODIR/mock/getter.go k8s.io/client-go/kubernetes/typed/rbac/v1 ClusterRolesGetter
//go:generate mockgen -copyright_file $CAPATH/LICENSE_HEADER -package=mock -destination=$GODIR/mock/interface.go k8s.io/client-go/kubernetes/typed/rbac/v1 ClusterRoleInterface
package cluster_role

import (
	rbac "k8s.io/api/rbac/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbac_type "k8s.io/client-go/kubernetes/typed/rbac/v1"

	"github.com/pkg/errors"

	cage_k8s "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1"
)

// Client provides an interface to cluster roles.
type Client interface {
	Get(role string, options ...meta.GetOptions) (_ *rbac.ClusterRole, exists bool, _ error)
	List(options ...meta.ListOptions) (*rbac.ClusterRoleList, error)
}

// DefaultClient implementation of Client operates on a real kubernetes API.
type DefaultClient struct {
	rbac_type.ClusterRolesGetter
}

// NewDefaultClient returns an initialized DefaultClient.
func NewDefaultClient(getter rbac_type.ClusterRolesGetter) *DefaultClient {
	return &DefaultClient{ClusterRolesGetter: getter}
}

// Get returns the object if found, reports that the object does not exist, or returns an error.
//
// A single GetOptions value can be passed as the final argument to customize the query.
func (c *DefaultClient) Get(role string, options ...meta.GetOptions) (_ *rbac.ClusterRole, exists bool, _ error) {
	obj, err := c.ClusterRoles().Get(role, cage_k8s.GetOptionsFromVariadic(options))
	if err != nil {
		if k8s_errors.IsNotFound(err) {
			return nil, false, nil
		}
		return nil, false, errors.Wrap(err, "failed to get cluster role")
	}

	return obj, true, nil
}

// List returns the matching objects.
//
// A single ListOptions value can be passed as the final argument to customize the query.
//
// It implements Client.
func (c *DefaultClient) List(options ...meta.ListOptions) (*rbac.ClusterRoleList, error) {
	list, err := c.ClusterRoles().List(cage_k8s.ListOptionsFromVariadic(options))
	if err != nil {
		if k8s_errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to list cluster roles")
	}

	return list, nil
}

var _ Client = (*DefaultClient)(nil)
