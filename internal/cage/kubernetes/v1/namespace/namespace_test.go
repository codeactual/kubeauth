// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package namespace_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/namespace"
	mock_namespace "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/namespace/mock"
	cage_k8s_testkit "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/testkit"
	cage_require "github.com/codeactual/kubeauth/internal/cage/testkit/testify/require"
)

const (
	Namespace = "some-namespace"
)

func newClient(mockCtrl *gomock.Controller) (*mock_namespace.MockNamespaceInterface, *namespace.DefaultClient) {
	mockInterface := mock_namespace.NewMockNamespaceInterface(mockCtrl)
	mockGetter := mock_namespace.NewMockNamespacesGetter(mockCtrl)
	mockGetter.EXPECT().Namespaces().Return(mockInterface)
	return mockInterface, namespace.NewDefaultClient(mockGetter)
}

func TestGet(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectOptions := meta.GetOptions{}
		expectNs := &core.Namespace{ObjectMeta: meta.ObjectMeta{Name: Namespace}}

		mockInterface, wrapperClient := newClient(mockCtrl)
		mockInterface.EXPECT().Get(Namespace, expectOptions).Return(expectNs, nil)

		actualNs, exists, err := wrapperClient.Get(Namespace, expectOptions)
		require.NoError(t, err)
		require.True(t, exists)
		require.Exactly(t, expectNs, actualNs)
	})

	t.Run("not found", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectOptions := meta.GetOptions{}

		mockInterface, wrapperClient := newClient(mockCtrl)
		mockInterface.EXPECT().Get(Namespace, expectOptions).Return(nil, cage_k8s_testkit.NotFound())

		actualNs, exists, err := wrapperClient.Get(Namespace, expectOptions)
		require.NoError(t, err)
		require.False(t, exists)
		require.Nil(t, actualNs)
	})

	t.Run("error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectOptions := meta.GetOptions{}
		expectErr := errors.New("expectErr")

		mockInterface, wrapperClient := newClient(mockCtrl)
		mockInterface.EXPECT().Get(Namespace, expectOptions).Return(nil, expectErr)

		actualNs, exists, actualErr := wrapperClient.Get(Namespace, expectOptions)
		cage_require.MatchRegexp(t, fmt.Sprintf("%v", actualErr), "failed to get namespace.*expectErr")
		require.False(t, exists)
		require.Nil(t, actualNs)
	})
}
