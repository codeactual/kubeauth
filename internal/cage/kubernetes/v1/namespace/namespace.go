// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

//go:generate mockgen -copyright_file=$LICENSE_HEADER -package=mock -destination=$GODIR/mock/wrapper.go -source=$GODIR/$GOFILE
//go:generate mockgen -copyright_file $CAPATH/LICENSE_HEADER -package=mock -destination=$GODIR/mock/getter.go k8s.io/client-go/kubernetes/typed/core/v1 NamespacesGetter
//go:generate mockgen -copyright_file $CAPATH/LICENSE_HEADER -package=mock -destination=$GODIR/mock/interface.go k8s.io/client-go/kubernetes/typed/core/v1 NamespaceInterface
package namespace

import (
	core "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	core_type "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/pkg/errors"

	cage_k8s "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1"
)

// Client provides an interface to namespaces.
type Client interface {
	Get(name string, options ...meta.GetOptions) (_ *core.Namespace, exists bool, _ error)
}

// DefaultClient implementation of Client operates on a real kubernetes API.
type DefaultClient struct {
	core_type.NamespacesGetter
}

// NewDefaultClient returns an initialized DefaultClient.
func NewDefaultClient(getter core_type.NamespacesGetter) *DefaultClient {
	return &DefaultClient{NamespacesGetter: getter}
}

// Get returns the object if found, reports that the object does not exist, or returns an error.
//
// A single GetOptions value can be passed as the final argument to customize the query.
func (c *DefaultClient) Get(name string, options ...meta.GetOptions) (_ *core.Namespace, exists bool, _ error) {
	obj, err := c.Namespaces().Get(name, cage_k8s.GetOptionsFromVariadic(options))
	if err != nil {
		if k8s_errors.IsNotFound(err) {
			return nil, false, nil
		}
		return nil, false, errors.Wrapf(err, "failed to get namespace [%s]", name)
	}

	return obj, true, nil
}
