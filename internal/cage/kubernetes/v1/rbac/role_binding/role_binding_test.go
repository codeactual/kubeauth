// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package role_binding_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	rbac "k8s.io/api/rbac/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	cage_k8s "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1"
	"github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/role_binding"
	mock_role_binding "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/role_binding/mock"
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

func newClient(mockCtrl *gomock.Controller, namespace string) (*mock_role_binding.MockRoleBindingInterface, *role_binding.DefaultClient) {
	mockInterface := mock_role_binding.NewMockRoleBindingInterface(mockCtrl)
	mockGetter := mock_role_binding.NewMockRoleBindingsGetter(mockCtrl)
	mockGetter.EXPECT().RoleBindings(namespace).Return(mockInterface)
	return mockInterface, role_binding.NewDefaultClient(mockGetter)
}

func TestCreate(t *testing.T) {
	t.Run("created", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectSubject := rbac.Subject{Name: SubjectName, Kind: SubjectKind, Namespace: SubjectNamespace}
		expectBinding := &rbac.RoleBinding{
			ObjectMeta: meta.ObjectMeta{Namespace: Namespace, Name: Binding},
			RoleRef:    rbac.RoleRef{Kind: cage_k8s.KindRole, Name: Role},
			Subjects:   []rbac.Subject{expectSubject},
		}

		mockInterface, wrapperClient := newClient(mockCtrl, Namespace)
		mockInterface.EXPECT().Create(expectBinding).Return(expectBinding, nil)

		actualBinding, err := wrapperClient.Create(Namespace, Binding, Role, expectSubject)
		require.NoError(t, err)
		require.Exactly(t, expectBinding, actualBinding)
	})

	t.Run("error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectSubject := rbac.Subject{Name: SubjectName, Kind: SubjectKind, Namespace: SubjectNamespace}
		expectBinding := &rbac.RoleBinding{
			ObjectMeta: meta.ObjectMeta{Namespace: Namespace, Name: Binding},
			RoleRef:    rbac.RoleRef{Kind: cage_k8s.KindRole, Name: Role},
			Subjects:   []rbac.Subject{expectSubject},
		}
		expectErr := errors.New("expectErr")

		mockInterface, wrapperClient := newClient(mockCtrl, Namespace)
		mockInterface.EXPECT().Create(expectBinding).Return(nil, expectErr)

		actualBinding, actualErr := wrapperClient.Create(Namespace, Binding, Role, expectSubject)
		cage_require.MatchRegexp(t, fmt.Sprintf("%v", actualErr), "failed to bind role.*expectErr")
		require.Nil(t, actualBinding)
	})
}

func TestList(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectOptions := meta.ListOptions{}
		expectSubject := rbac.Subject{Name: SubjectName, Kind: SubjectKind, Namespace: SubjectNamespace}
		expectList := &rbac.RoleBindingList{
			Items: []rbac.RoleBinding{
				{
					ObjectMeta: meta.ObjectMeta{Name: Binding},
					RoleRef:    rbac.RoleRef{Kind: cage_k8s.KindRole, Name: Role},
					Subjects:   []rbac.Subject{expectSubject},
				},
			},
		}

		mockInterface, wrapperClient := newClient(mockCtrl, Namespace)
		mockInterface.EXPECT().List(expectOptions).Return(expectList, nil)

		actualList, actualErr := wrapperClient.List(Namespace, expectOptions)
		require.NoError(t, actualErr)
		require.Exactly(t, expectList, actualList)
	})

	t.Run("error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectOptions := meta.ListOptions{}
		expectErr := errors.New("expectErr")

		mockInterface, wrapperClient := newClient(mockCtrl, Namespace)
		mockInterface.EXPECT().List(expectOptions).Return(nil, expectErr)

		actualRole, actualErr := wrapperClient.List(Namespace, expectOptions)
		cage_require.MatchRegexp(t, fmt.Sprintf("%v", actualErr), "failed to list role binding.*expectErr")
		require.Nil(t, actualRole)
	})
}
