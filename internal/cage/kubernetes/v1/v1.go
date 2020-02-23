// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package v1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	EmptyNamespace = "<no namespace>"

	KindClusterRole        = "ClusterRole"
	KindClusterRoleBinding = "ClusterRoleBinding"
	KindGroup              = "Group"
	KindRole               = "Role"
	KindRoleBinding        = "RoleBinding"
	KindServiceAccount     = "ServiceAccount"
	KindUser               = "User"
)

// GetOptionsFromVariadic returns the first input GetOptions element or a zero-value GetOptions.
func GetOptionsFromVariadic(variadic []meta.GetOptions) meta.GetOptions {
	if len(variadic) == 0 {
		variadic = append(variadic, meta.GetOptions{})
	}
	return variadic[0]
}

// ListOptionsFromVariadic returns the first input ListOptions element or a zero-value ListOptions.
func ListOptionsFromVariadic(variadic []meta.ListOptions) meta.ListOptions {
	if len(variadic) == 0 {
		variadic = append(variadic, meta.ListOptions{})
	}
	return variadic[0]
}
