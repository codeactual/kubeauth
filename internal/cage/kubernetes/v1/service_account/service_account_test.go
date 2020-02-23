// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package service_account_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/service_account"
	mock_sa "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/service_account/mock"
	cage_k8s_testkit "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/testkit"
	cage_require "github.com/codeactual/kubeauth/internal/cage/testkit/testify/require"
)

const (
	Namespace      = "some-namespace"
	ServiceAccount = "some-sa"
)

func newClient(mockCtrl *gomock.Controller, namespace string) (*mock_sa.MockServiceAccountInterface, *service_account.DefaultClient) {
	mockInterface := mock_sa.NewMockServiceAccountInterface(mockCtrl)
	mockGetter := mock_sa.NewMockServiceAccountsGetter(mockCtrl)
	mockGetter.EXPECT().ServiceAccounts(namespace).Return(mockInterface)
	return mockInterface, service_account.NewDefaultClient(mockGetter)
}

func TestCreateBasic(t *testing.T) {
	t.Run("created", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectSa := &core.ServiceAccount{ObjectMeta: meta.ObjectMeta{Name: ServiceAccount}}

		mockInterface, wrapperClient := newClient(mockCtrl, Namespace)
		mockInterface.EXPECT().Create(expectSa).Return(expectSa, nil)

		actualSa, err := wrapperClient.CreateBasic(Namespace, ServiceAccount)
		require.NoError(t, err)
		require.Exactly(t, expectSa, actualSa)
	})

	t.Run("error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectSa := &core.ServiceAccount{ObjectMeta: meta.ObjectMeta{Name: ServiceAccount}}
		expectErr := errors.New("expectErr")

		mockInterface, wrapperClient := newClient(mockCtrl, Namespace)
		mockInterface.EXPECT().Create(expectSa).Return(nil, expectErr)

		actualSa, actualErr := wrapperClient.CreateBasic(Namespace, ServiceAccount)
		cage_require.MatchRegexp(t, fmt.Sprintf("%v", actualErr), "failed to create service account.*expectErr")
		require.Nil(t, actualSa)
	})
}

func TestGet(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectOptions := meta.GetOptions{}
		expectSa := &core.ServiceAccount{ObjectMeta: meta.ObjectMeta{Name: ServiceAccount}}

		mockInterface, wrapperClient := newClient(mockCtrl, Namespace)
		mockInterface.EXPECT().Get(ServiceAccount, expectOptions).Return(expectSa, nil)

		actualSa, exists, err := wrapperClient.Get(Namespace, ServiceAccount, expectOptions)
		require.NoError(t, err)
		require.True(t, exists)
		require.Exactly(t, expectSa, actualSa)
	})

	t.Run("not found", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectOptions := meta.GetOptions{}

		mockInterface, wrapperClient := newClient(mockCtrl, Namespace)
		mockInterface.EXPECT().Get(ServiceAccount, expectOptions).Return(nil, cage_k8s_testkit.NotFound())

		actualSa, exists, err := wrapperClient.Get(Namespace, ServiceAccount, expectOptions)
		require.NoError(t, err)
		require.False(t, exists)
		require.Nil(t, actualSa)
	})

	t.Run("error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectOptions := meta.GetOptions{}
		expectErr := errors.New("expectErr")

		mockInterface, wrapperClient := newClient(mockCtrl, Namespace)
		mockInterface.EXPECT().Get(ServiceAccount, expectOptions).Return(nil, expectErr)

		actualSa, exists, actualErr := wrapperClient.Get(Namespace, ServiceAccount, expectOptions)
		cage_require.MatchRegexp(t, fmt.Sprintf("%v", actualErr), "failed to get service account.*expectErr")
		require.False(t, exists)
		require.Nil(t, actualSa)
	})
}

func TestList(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		expectOptions := meta.ListOptions{}
		expectList := &core.ServiceAccountList{
			Items: []core.ServiceAccount{
				{ObjectMeta: meta.ObjectMeta{Name: ServiceAccount}},
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

		actualList, actualErr := wrapperClient.List(Namespace, expectOptions)
		cage_require.MatchRegexp(t, fmt.Sprintf("%v", actualErr), "failed to list service accounts.*expectErr")
		require.Nil(t, actualList)
	})
}
