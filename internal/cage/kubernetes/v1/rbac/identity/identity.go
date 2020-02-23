// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package identity

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IdentitySource describes where an Idenity was found, e.g. RoleBinding.
type IdentitySource struct {
	meta.TypeMeta
	meta.ObjectMeta
}

// String returns the relevant fields in a human-readable format for use in info/error messages.
func (i IdentitySource) String() (s string) {
	kind := i.Kind

	// As if v0.16.4, some list result objects contain an empty Kind value. For processors which do not
	// provide a default value, ensure we include one here.
	if kind == "" {
		kind = "<empty Kind>"
	}

	s = kind + " " + i.Name
	if i.Namespace != "" {
		s += " of namespace " + i.Namespace
	}
	return s
}

// Identity describes an object which may have RBAC grants.
type Identity struct {
	meta.TypeMeta
	meta.ObjectMeta

	// Source describes the object (if any) in which this Identity was found, e.g. RoleBinding.
	Source *IdentitySource

	// Querier indicates which IdentityQuerier implementation produced this value.
	Querier string
}

// String returns the relevant fields in a human-readable format for use in info/error messages.
func (i Identity) String() (s string) {
	kind := i.Kind

	// As if v0.16.4, some list result objects contain an empty Kind value. For processors which do not
	// provide a default value, ensure we include one here.
	if kind == "" {
		kind = "<empty Kind>"
	}

	s = kind + " " + i.Name
	if i.Namespace != "" {
		s += " of namespace " + i.Namespace
	}
	if i.Source != nil {
		s += " (from " + i.Source.String() + ")"
	}
	s += " via [" + i.Querier + "] querier"
	return s
}

// IdentityList is a collection of Identity values.
//
// Its structure ("Items") aligns with the list collections in k8s.io/api/rbac/v1.
type IdentityList struct {
	// Items holds the collection elements.
	Items []Identity
}

// Add appends and returns a new list item.
func (i *IdentityList) Add(namespace, kind, name string, source *IdentitySource) {
	i.Items = append(i.Items, Identity{
		ObjectMeta: meta.ObjectMeta{Namespace: namespace, Name: name},
		TypeMeta:   meta.TypeMeta{Kind: kind},
		Source:     source,
	})
}
