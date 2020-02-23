// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cluster_role_binding_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	rbac "k8s.io/api/rbac/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	cage_k8s "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1"
	"github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/cluster_role_binding"
	mock_cluster_role_binding "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/cluster_role_binding/mock"
	cage_require "github.com/codeactual/kubeauth/internal/cage/testkit/testify/require"
)

const (
	Namespace        = "some-namespace"
	Role             = "some-role"
	Binding          = "some-binding"
	SubjectKind      = "ServiceAccount"
	SubjectName      = "some-sa"
	SubjectNamespace = "some-sa-namespace"
)

func newClient(mockCtrl *gomock.Controller) (*mock_cluster_role_binding.MockClusterRoleBindingInterface, *cluster_role_binding.DefaultClient) {
	mockInterface := mock_cluster_role_binding.NewMockClusterRoleBindingInterface(mockCtrl)
	mockGetter := mock_cluster_role_binding.NewMockClusterRoleBindingsGetter(mockCtrl)
	mockGetter.EXPECT().ClusterRoleBindings().Return(mockInterface)
	return mockInterface, cluster_role_binding.NewDefaultClient(mockGetter)
}

func TestCreate(t *testing.T) {
	t.Run("created", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectSubject := rbac.Subject{Name: SubjectName, Kind: SubjectKind, Namespace: SubjectNamespace}
		expectBinding := &rbac.ClusterRoleBinding{
			ObjectMeta: meta.ObjectMeta{Name: Binding},
			RoleRef:    rbac.RoleRef{Kind: cage_k8s.KindClusterRole, Name: Role},
			Subjects:   []rbac.Subject{expectSubject},
		}

		mockInterface, wrapperClient := newClient(mockCtrl)
		mockInterface.EXPECT().Create(expectBinding).Return(expectBinding, nil)

		actualBinding, err := wrapperClient.Create(Binding, Role, expectSubject)
		require.NoError(t, err)
		require.Exactly(t, expectBinding, actualBinding)
	})

	t.Run("error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectSubject := rbac.Subject{Name: SubjectName, Kind: SubjectKind, Namespace: SubjectNamespace}
		expectBinding := &rbac.ClusterRoleBinding{
			ObjectMeta: meta.ObjectMeta{Name: Binding},
			RoleRef:    rbac.RoleRef{Kind: cage_k8s.KindClusterRole, Name: Role},
			Subjects:   []rbac.Subject{expectSubject},
		}
		expectErr := errors.New("expectErr")

		mockInterface, wrapperClient := newClient(mockCtrl)
		mockInterface.EXPECT().Create(expectBinding).Return(nil, expectErr)

		actualBinding, actualErr := wrapperClient.Create(Binding, Role, expectSubject)
		cage_require.MatchRegexp(t, fmt.Sprintf("%v", actualErr), "failed to bind cluster role.*expectErr")
		require.Nil(t, actualBinding)
	})
}

func TestList(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectOptions := meta.ListOptions{}
		expectSubject := rbac.Subject{Name: SubjectName, Kind: SubjectKind, Namespace: SubjectNamespace}
		expectList := &rbac.ClusterRoleBindingList{
			Items: []rbac.ClusterRoleBinding{
				{
					ObjectMeta: meta.ObjectMeta{Name: Binding},
					RoleRef:    rbac.RoleRef{Kind: cage_k8s.KindClusterRole, Name: Role},
					Subjects:   []rbac.Subject{expectSubject},
				},
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

		actualRole, actualErr := wrapperClient.List(expectOptions)
		cage_require.MatchRegexp(t, fmt.Sprintf("%v", actualErr), "failed to list cluster role binding.*expectErr")
		require.Nil(t, actualRole)
	})
}
