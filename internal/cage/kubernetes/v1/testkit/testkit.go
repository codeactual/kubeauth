// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package testkit

import (
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// NotFound returns a resource-agnostic error that satisfies k8s.io/apimachinery/pkg/api/errors.IsNotFound.
func NotFound() error {
	return k8s_errors.NewNotFound(schema.GroupResource{}, "")
}
