// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cluster_role_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	rbac "k8s.io/api/rbac/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/cluster_role"
	mock_cluster_role "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/cluster_role/mock"
	cage_k8s_testkit "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/testkit"
	cage_require "github.com/codeactual/kubeauth/internal/cage/testkit/testify/require"
)

const (
	Namespace   = "some-namespace"
	ClusterRole = "some-cluster-role"
)

func newClient(mockCtrl *gomock.Controller) (*mock_cluster_role.MockClusterRoleInterface, *cluster_role.DefaultClient) {
	mockInterface := mock_cluster_role.NewMockClusterRoleInterface(mockCtrl)
	mockGetter := mock_cluster_role.NewMockClusterRolesGetter(mockCtrl)
	mockGetter.EXPECT().ClusterRoles().Return(mockInterface)
	return mockInterface, cluster_role.NewDefaultClient(mockGetter)
}

func TestGet(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectOptions := meta.GetOptions{}
		expectRole := &rbac.ClusterRole{ObjectMeta: meta.ObjectMeta{Name: ClusterRole}}

		mockInterface, wrapperClient := newClient(mockCtrl)
		mockInterface.EXPECT().Get(ClusterRole, expectOptions).Return(expectRole, nil)

		actualRole, exists, err := wrapperClient.Get(ClusterRole, expectOptions)
		require.NoError(t, err)
		require.True(t, exists)
		require.Exactly(t, expectRole, actualRole)
	})

	t.Run("not found", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectOptions := meta.GetOptions{}

		mockInterface, wrapperClient := newClient(mockCtrl)
		mockInterface.EXPECT().Get(ClusterRole, expectOptions).Return(nil, cage_k8s_testkit.NotFound())

		actualRole, exists, err := wrapperClient.Get(ClusterRole, expectOptions)
		require.NoError(t, err)
		require.False(t, exists)
		require.Nil(t, actualRole)
	})

	t.Run("error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectOptions := meta.GetOptions{}
		expectErr := errors.New("expectErr")

		mockInterface, wrapperClient := newClient(mockCtrl)
		mockInterface.EXPECT().Get(ClusterRole, expectOptions).Return(nil, expectErr)

		actualRole, exists, actualErr := wrapperClient.Get(ClusterRole, expectOptions)
		cage_require.MatchRegexp(t, fmt.Sprintf("%v", actualErr), "failed to get cluster role.*expectErr")
		require.False(t, exists)
		require.Nil(t, actualRole)
	})
}

func TestList(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectOptions := meta.ListOptions{}
		expectList := &rbac.ClusterRoleList{
			Items: []rbac.ClusterRole{
				{ObjectMeta: meta.ObjectMeta{Name: ClusterRole}},
			},
		}

		mockInterface, wrapperClient := newClient(mockCtrl)
		mockInterface.EXPECT().List(expectOptions).Return(expectList, nil)

		actualList, actualErr := wrapperClient.List(expectOptions)
		require.NoError(t, actualErr)
		require.Exactly(t, expectList, actualList)
	})

	t.Run("error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectOptions := meta.ListOptions{}
		expectErr := errors.New("expectErr")

		mockInterface, wrapperClient := newClient(mockCtrl)
		mockInterface.EXPECT().List(expectOptions).Return(nil, expectErr)

		actualList, actualErr := wrapperClient.List(expectOptions)
		cage_require.MatchRegexp(t, fmt.Sprintf("%v", actualErr), "failed to list cluster role.*expectErr")
		require.Nil(t, actualList)
	})
}
