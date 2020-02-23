// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

//go:generate mockgen -copyright_file=$LICENSE_HEADER -package=mock -destination=$GODIR/mock/wrapper.go -source=$GODIR/$GOFILE
//go:generate mockgen -copyright_file $CAPATH/LICENSE_HEADER -package=mock -destination=$GODIR/mock/getter.go k8s.io/client-go/kubernetes/typed/core/v1 ServiceAccountsGetter
//go:generate mockgen -copyright_file $CAPATH/LICENSE_HEADER -package=mock -destination=$GODIR/mock/interface.go k8s.io/client-go/kubernetes/typed/core/v1 ServiceAccountInterface
package service_account

import (
	core "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	core_type "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/pkg/errors"

	cage_k8s "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1"
)

// Client provides an interface to service accounts.
type Client interface {
	CreateBasic(ns, sa string) (*core.ServiceAccount, error)
	Get(ns, sa string, options ...meta.GetOptions) (_ *core.ServiceAccount, exists bool, _ error)
	List(ns string, options ...meta.ListOptions) (*core.ServiceAccountList, error)
}

// DefaultClient implementation of Client operates on a real kubernetes API.
type DefaultClient struct {
	core_type.ServiceAccountsGetter
}

// NewDefaultClient returns an initialized DefaultClient.
func NewDefaultClient(getter core_type.ServiceAccountsGetter) *DefaultClient {
	return &DefaultClient{ServiceAccountsGetter: getter}
}

// CreateBasic adds a service account based on only its namespace and account name.
func (c *DefaultClient) CreateBasic(ns, sa string) (*core.ServiceAccount, error) {
	created, err := c.ServiceAccounts(ns).Create(&core.ServiceAccount{ObjectMeta: meta.ObjectMeta{Name: sa}})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create service account [%s] in namespace [%s]", sa, ns)
	}

	return created, nil
}

// Get returns the object if found, reports that the object does not exist, or returns an error.
//
// A single GetOptions value can be passed as the final argument to customize the query.
func (c *DefaultClient) Get(ns, sa string, options ...meta.GetOptions) (_ *core.ServiceAccount, exists bool, _ error) {
	obj, err := c.ServiceAccounts(ns).Get(sa, cage_k8s.GetOptionsFromVariadic(options))
	if err != nil {
		if k8s_errors.IsNotFound(err) {
			return nil, false, nil
		}
		return nil, false, errors.Wrapf(err, "failed to get service account [%s] in namespace [%s]", sa, ns)
	}

	return obj, true, nil
}

// List returns the matching objects.
//
// A single ListOptions value can be passed as the final argument to customize the query.
//
// It implements Client.
func (c *DefaultClient) List(ns string, options ...meta.ListOptions) (*core.ServiceAccountList, error) {
	list, err := c.ServiceAccounts(ns).List(cage_k8s.ListOptionsFromVariadic(options))
	if err != nil {
		if k8s_errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to list service accounts")
	}

	return list, nil
}

var _ Client = (*DefaultClient)(nil)
