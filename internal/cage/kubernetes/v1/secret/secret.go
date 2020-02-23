// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

//go:generate mockgen -copyright_file=$LICENSE_HEADER -package=mock -destination=$GODIR/mock/wrapper.go -source=$GODIR/$GOFILE
//go:generate mockgen -copyright_file $CAPATH/LICENSE_HEADER -package=mock -destination=$GODIR/mock/getter.go k8s.io/client-go/kubernetes/typed/core/v1 SecretsGetter
//go:generate mockgen -copyright_file $CAPATH/LICENSE_HEADER -package=mock -destination=$GODIR/mock/interface.go k8s.io/client-go/kubernetes/typed/core/v1 SecretInterface
package secret

import (
	core "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	core_type "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/pkg/errors"

	cage_k8s "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1"
)

// Client provides an interface to secrets.
type Client interface {
	Get(ns, name string, options ...meta.GetOptions) (_ *core.Secret, exists bool, _ error)
}

// DefaultClient implementation of Client operates on a real kubernetes API.
type DefaultClient struct {
	core_type.SecretsGetter
}

// NewDefaultClient returns an initialized DefaultClient.
func NewDefaultClient(getter core_type.SecretsGetter) *DefaultClient {
	return &DefaultClient{SecretsGetter: getter}
}

// Get returns the secret object if found, reports that the object does not exist,
// or returns an error.
//
// A single GetOptions value can be passed as the final argument to customize the query.
//
// It implements Client.
func (c *DefaultClient) Get(ns, name string, options ...meta.GetOptions) (_ *core.Secret, exists bool, _ error) {
	obj, err := c.Secrets(ns).Get(name, cage_k8s.GetOptionsFromVariadic(options))
	if err != nil {
		if k8s_errors.IsNotFound(err) {
			return nil, false, nil
		}
		return nil, false, errors.Wrapf(err, "failed to get secret [%s] in namespace [%s]", name, ns)
	}

	return obj, true, nil
}

var _ Client = (*DefaultClient)(nil)
