// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package identity_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	cage_k8s "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1"
	cage_k8s_core "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/core"
	mock_core "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/core/mock"
	cage_k8s_identity "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/identity"
)

const (
	DoesNotExist = "does-not-exist"

	Exists           = true
	NotExists        = false
	NoQueryNamespace = ""

	ConfigUsername                      = "some-username"
	ContextName                         = "some-context-name"
	Namespace                           = "some-namespace"
	ServiceAccountUsernameBase          = "some-sa"
	ServiceAccountUsername              = "system:serviceaccount:" + Namespace + ":" + ServiceAccountUsernameBase
	ServiceAccountUsernameMiss          = "system:serviceaccount:" + Namespace + ":" + DoesNotExist
	ServiceAccountUsernameNamespaceMiss = "system:serviceaccount:" + DoesNotExist + ":" + ServiceAccountUsernameBase
	ServiceAccountGroupOneNamespace     = "system:serviceaccounts:" + Namespace
	ServiceAccountGroupAllNamespace     = "system:serviceaccounts"
	ServiceAccountGroupNamespaceMiss    = "system:serviceaccounts:" + DoesNotExist

	CoreUsername = "system:anonymous"
	CoreGroup    = "system:masters"
)

func ctx() context.Context {
	return context.Background()
}

func newClientCmdConfigWithUser(namespace, username string) *clientcmdapi.Config {
	return &clientcmdapi.Config{
		Contexts: map[string]*clientcmdapi.Context{
			ContextName: {
				AuthInfo:  username,
				Namespace: namespace,
			},
		},
	}
}

func TestConfigUserQuerier(t *testing.T) {
	// Assert that if the query does not require a specific object kind, still support the query
	// and filter on name/namespace.
	t.Run("compatible with empty kind", func(t *testing.T) {
		require.True(t, cage_k8s_identity.ConfigUserQuerier{}.Compatible(&cage_k8s_identity.Query{}))
	})

	t.Run("incompatible with non user kind", func(t *testing.T) {
		require.False(t, cage_k8s_identity.ConfigUserQuerier{}.Compatible(&cage_k8s_identity.Query{Kind: cage_k8s.KindGroup}))
	})

	t.Run("hit", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind:            cage_k8s.KindUser,
			Name:            ConfigUsername,
			ClientCmdConfig: newClientCmdConfigWithUser(Namespace, ConfigUsername),
		}
		var nonSut *cage_k8s_core.Clientset

		list, err := cage_k8s_identity.ConfigUserQuerier{}.Do(ctx(), nonSut, &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		require.Exactly(t, cage_k8s.KindUser, list.Items[0].Kind)
		require.Exactly(t, ConfigUsername, list.Items[0].Name)
		require.Exactly(t, Namespace, list.Items[0].Namespace)
	})

	t.Run("name miss", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind:            cage_k8s.KindUser,
			Name:            DoesNotExist,
			ClientCmdConfig: newClientCmdConfigWithUser(Namespace, ConfigUsername),
		}
		var nonSut *cage_k8s_core.Clientset

		list, err := cage_k8s_identity.ConfigUserQuerier{}.Do(ctx(), nonSut, &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 0)
	})

	t.Run("namespace miss", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind:            cage_k8s.KindUser,
			Name:            ConfigUsername,
			Namespace:       Namespace,
			ClientCmdConfig: newClientCmdConfigWithUser(DoesNotExist, ConfigUsername),
		}
		var nonSut *cage_k8s_core.Clientset

		list, err := cage_k8s_identity.ConfigUserQuerier{}.Do(ctx(), nonSut, &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 0)
	})
}

func TestCoreUserQuerier(t *testing.T) {
	// Assert that if the query does not require a specific object kind, still support the query
	// and filter on name/namespace.
	t.Run("compatible with empty kind", func(t *testing.T) {
		require.True(t, cage_k8s_identity.CoreUserQuerier{}.Compatible(&cage_k8s_identity.Query{}))
	})

	t.Run("incompatible with non user kind", func(t *testing.T) {
		require.False(t, cage_k8s_identity.CoreUserQuerier{}.Compatible(&cage_k8s_identity.Query{Kind: cage_k8s.KindGroup}))
	})

	t.Run("hit", func(t *testing.T) {
		username := CoreUsername
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindUser,
			Name: username,
		}
		var nonSut *cage_k8s_core.Clientset

		list, err := cage_k8s_identity.CoreUserQuerier{}.Do(ctx(), nonSut, &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		require.Exactly(t, cage_k8s.KindUser, list.Items[0].Kind)
		require.Exactly(t, username, list.Items[0].Name)
	})

	t.Run("name miss", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindUser,
			Name: DoesNotExist,
		}
		var nonSut *cage_k8s_core.Clientset

		list, err := cage_k8s_identity.CoreUserQuerier{}.Do(ctx(), nonSut, &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 0)
	})
}

func TestCoreGroupQuerier(t *testing.T) {
	// Assert that if the query does not require a specific object kind, still support the query
	// and filter on name/namespace.
	t.Run("compatible with empty kind", func(t *testing.T) {
		require.True(t, cage_k8s_identity.CoreGroupQuerier{}.Compatible(&cage_k8s_identity.Query{}))
	})

	t.Run("incompatible with non group kind", func(t *testing.T) {
		require.False(t, cage_k8s_identity.CoreGroupQuerier{}.Compatible(&cage_k8s_identity.Query{Kind: cage_k8s.KindUser}))
	})

	t.Run("hit", func(t *testing.T) {
		group := CoreGroup
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindGroup,
			Name: group,
		}
		var nonSut *cage_k8s_core.Clientset

		list, err := cage_k8s_identity.CoreGroupQuerier{}.Do(ctx(), nonSut, &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		require.Exactly(t, cage_k8s.KindGroup, list.Items[0].Kind)
		require.Exactly(t, group, list.Items[0].Name)
	})

	t.Run("name miss", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindUser,
			Name: DoesNotExist,
		}
		var nonSut *cage_k8s_core.Clientset

		list, err := cage_k8s_identity.CoreGroupQuerier{}.Do(ctx(), nonSut, &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 0)
	})
}

func TestRoleSubjectQuerier(t *testing.T) {
	// Assert that if the query does not require a specific object kind, still support the query
	// and filter on name/namespace.
	t.Run("compatible with empty kind", func(t *testing.T) {
		require.True(t, cage_k8s_identity.RoleSubjectQuerier{}.Compatible(&cage_k8s_identity.Query{}))
	})

	t.Run("compatible with subject kind", func(t *testing.T) {
		require.True(t, cage_k8s_identity.RoleSubjectQuerier{}.Compatible(&cage_k8s_identity.Query{Kind: cage_k8s.KindUser}))
		require.True(t, cage_k8s_identity.RoleSubjectQuerier{}.Compatible(&cage_k8s_identity.Query{Kind: cage_k8s.KindGroup}))
	})

	t.Run("incompatible with non subject kind", func(t *testing.T) {
		require.False(t, cage_k8s_identity.RoleSubjectQuerier{}.Compatible(&cage_k8s_identity.Query{Kind: cage_k8s.KindRole}))
	})

	t.Run("user kind hit", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindUser,
			Name: CoreUsername,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		bindings := rbac.RoleBindingList{
			Items: []rbac.RoleBinding{
				{
					Subjects: []rbac.Subject{
						{Kind: cage_k8s.KindUser, Name: CoreUsername, Namespace: Namespace},
					},
				},
			},
		}
		mockClientset.RoleBindings.EXPECT().List(NoQueryNamespace, meta.ListOptions{}).Return(&bindings, nil)

		list, err := cage_k8s_identity.RoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		require.Exactly(t, cage_k8s.KindUser, list.Items[0].Kind)
		require.Exactly(t, CoreUsername, list.Items[0].Name)
	})

	t.Run("group kind hit", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindGroup,
			Name: CoreGroup,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		bindings := rbac.RoleBindingList{
			Items: []rbac.RoleBinding{
				{
					Subjects: []rbac.Subject{
						{Kind: cage_k8s.KindGroup, Name: CoreGroup, Namespace: Namespace},
					},
				},
			},
		}
		mockClientset.RoleBindings.EXPECT().List(NoQueryNamespace, meta.ListOptions{}).Return(&bindings, nil)

		list, err := cage_k8s_identity.RoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		require.Exactly(t, cage_k8s.KindGroup, list.Items[0].Kind)
		require.Exactly(t, CoreGroup, list.Items[0].Name)
	})

	t.Run("service account user hit", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindServiceAccount,
			Name: ServiceAccountUsername,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		var nonSut *core.Namespace
		mockClientset.Namespaces.EXPECT().Get(Namespace).Return(nonSut, Exists, nil)

		bindings := rbac.RoleBindingList{
			Items: []rbac.RoleBinding{
				{
					Subjects: []rbac.Subject{
						{Kind: cage_k8s.KindServiceAccount, Name: ServiceAccountUsernameBase, Namespace: Namespace},
					},
				},
			},
		}
		mockClientset.RoleBindings.EXPECT().List(NoQueryNamespace, meta.ListOptions{}).Return(&bindings, nil)

		list, err := cage_k8s_identity.RoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		require.Exactly(t, cage_k8s.KindServiceAccount, list.Items[0].Kind)
		require.Exactly(t, ServiceAccountUsernameBase, list.Items[0].Name)
	})

	t.Run("service account one namespace group hit", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindGroup,
			Name: ServiceAccountGroupOneNamespace,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		var nonSut *core.Namespace
		mockClientset.Namespaces.EXPECT().Get(Namespace).Return(nonSut, Exists, nil)

		bindings := rbac.RoleBindingList{
			Items: []rbac.RoleBinding{
				{
					Subjects: []rbac.Subject{
						{Kind: cage_k8s.KindGroup, Name: ServiceAccountGroupOneNamespace, Namespace: Namespace},
					},
				},
			},
		}
		mockClientset.RoleBindings.EXPECT().List(NoQueryNamespace, meta.ListOptions{}).Return(&bindings, nil)

		list, err := cage_k8s_identity.RoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		require.Exactly(t, cage_k8s.KindGroup, list.Items[0].Kind)
		require.Exactly(t, ServiceAccountGroupOneNamespace, list.Items[0].Name)
	})

	t.Run("service account all namespace group hit", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindGroup,
			Name: ServiceAccountGroupAllNamespace,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		bindings := rbac.RoleBindingList{
			Items: []rbac.RoleBinding{
				{
					Subjects: []rbac.Subject{
						{Kind: cage_k8s.KindGroup, Name: ServiceAccountGroupAllNamespace, Namespace: Namespace},
					},
				},
			},
		}
		mockClientset.RoleBindings.EXPECT().List(NoQueryNamespace, meta.ListOptions{}).Return(&bindings, nil)

		list, err := cage_k8s_identity.RoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		require.Exactly(t, cage_k8s.KindGroup, list.Items[0].Kind)
		require.Exactly(t, ServiceAccountGroupAllNamespace, list.Items[0].Name)
	})

	t.Run("name miss", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindUser,
			Name: DoesNotExist,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		bindings := rbac.RoleBindingList{Items: []rbac.RoleBinding{}}
		mockClientset.RoleBindings.EXPECT().List(NoQueryNamespace, meta.ListOptions{}).Return(&bindings, nil)

		list, err := cage_k8s_identity.RoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 0)
	})

	t.Run("namespace miss", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind:      cage_k8s.KindUser,
			Name:      ServiceAccountUsername,
			Namespace: Namespace,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		var nonSut *core.Namespace
		mockClientset.Namespaces.EXPECT().Get(Namespace).Return(nonSut, Exists, nil)

		bindings := rbac.RoleBindingList{Items: []rbac.RoleBinding{}}
		mockClientset.RoleBindings.EXPECT().List(Namespace, meta.ListOptions{}).Return(&bindings, nil)

		list, err := cage_k8s_identity.RoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 0)
	})

	t.Run("service account namespace unknown", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindUser,
			Name: ServiceAccountUsername,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		var nonSut *core.Namespace
		mockClientset.Namespaces.EXPECT().Get(Namespace).Return(nonSut, NotExists, nil)

		list, err := cage_k8s_identity.RoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.Nil(t, list)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("service account [%s] namespace [%s] not found", ServiceAccountUsername, Namespace))
	})

	// Assert that an error is returned if a service account's namespace does not match the query's.
	t.Run("service account namespace mismatch", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind:      cage_k8s.KindUser,
			Name:      ServiceAccountUsername,
			Namespace: DoesNotExist,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		list, err := cage_k8s_identity.RoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.Nil(t, list)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("query's namespace [%s] does not match query service account [%s]'s namespace [%s] ", DoesNotExist, ServiceAccountUsername, Namespace))
	})
}

func TestClusterRoleSubjectQuerier(t *testing.T) {
	// Assert that if the query does not require a specific object kind, still support the query
	// and filter on name/namespace.
	t.Run("compatible with empty kind", func(t *testing.T) {
		require.True(t, cage_k8s_identity.ClusterRoleSubjectQuerier{}.Compatible(&cage_k8s_identity.Query{}))
	})

	t.Run("compatible with subject kind", func(t *testing.T) {
		require.True(t, cage_k8s_identity.ClusterRoleSubjectQuerier{}.Compatible(&cage_k8s_identity.Query{Kind: cage_k8s.KindUser}))
		require.True(t, cage_k8s_identity.ClusterRoleSubjectQuerier{}.Compatible(&cage_k8s_identity.Query{Kind: cage_k8s.KindGroup}))
	})

	t.Run("incompatible with non subject kind", func(t *testing.T) {
		require.False(t, cage_k8s_identity.ClusterRoleSubjectQuerier{}.Compatible(&cage_k8s_identity.Query{Kind: cage_k8s.KindRole}))
	})

	t.Run("user kind hit", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindUser,
			Name: CoreUsername,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		bindings := rbac.ClusterRoleBindingList{
			Items: []rbac.ClusterRoleBinding{
				{
					Subjects: []rbac.Subject{
						{Kind: cage_k8s.KindUser, Name: CoreUsername, Namespace: Namespace},
					},
				},
			},
		}
		mockClientset.ClusterRoleBindings.EXPECT().List(meta.ListOptions{}).Return(&bindings, nil)

		list, err := cage_k8s_identity.ClusterRoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		require.Exactly(t, cage_k8s.KindUser, list.Items[0].Kind)
		require.Exactly(t, CoreUsername, list.Items[0].Name)
	})

	t.Run("group kind hit", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindGroup,
			Name: CoreGroup,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		bindings := rbac.ClusterRoleBindingList{
			Items: []rbac.ClusterRoleBinding{
				{
					Subjects: []rbac.Subject{
						{Kind: cage_k8s.KindGroup, Name: CoreGroup, Namespace: Namespace},
					},
				},
			},
		}
		mockClientset.ClusterRoleBindings.EXPECT().List(meta.ListOptions{}).Return(&bindings, nil)

		list, err := cage_k8s_identity.ClusterRoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		require.Exactly(t, cage_k8s.KindGroup, list.Items[0].Kind)
		require.Exactly(t, CoreGroup, list.Items[0].Name)
	})

	t.Run("service account kind hit", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindUser,
			Name: ServiceAccountUsername,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		var nonSut *core.Namespace
		mockClientset.Namespaces.EXPECT().Get(Namespace).Return(nonSut, Exists, nil)

		bindings := rbac.ClusterRoleBindingList{
			Items: []rbac.ClusterRoleBinding{
				{
					Subjects: []rbac.Subject{
						{Kind: cage_k8s.KindServiceAccount, Name: ServiceAccountUsernameBase, Namespace: Namespace},
					},
				},
			},
		}
		mockClientset.ClusterRoleBindings.EXPECT().List(meta.ListOptions{}).Return(&bindings, nil)

		list, err := cage_k8s_identity.ClusterRoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		require.Exactly(t, cage_k8s.KindServiceAccount, list.Items[0].Kind)
		require.Exactly(t, ServiceAccountUsernameBase, list.Items[0].Name)
	})

	t.Run("service account one namespace group hit", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindGroup,
			Name: ServiceAccountGroupOneNamespace,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		var nonSut *core.Namespace
		mockClientset.Namespaces.EXPECT().Get(Namespace).Return(nonSut, Exists, nil)

		bindings := rbac.ClusterRoleBindingList{
			Items: []rbac.ClusterRoleBinding{
				{
					Subjects: []rbac.Subject{
						{Kind: cage_k8s.KindGroup, Name: ServiceAccountGroupOneNamespace, Namespace: Namespace},
					},
				},
			},
		}
		mockClientset.ClusterRoleBindings.EXPECT().List(meta.ListOptions{}).Return(&bindings, nil)

		list, err := cage_k8s_identity.ClusterRoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		require.Exactly(t, cage_k8s.KindGroup, list.Items[0].Kind)
		require.Exactly(t, ServiceAccountGroupOneNamespace, list.Items[0].Name)
	})

	t.Run("service account all namespace group hit", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindGroup,
			Name: ServiceAccountGroupAllNamespace,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		bindings := rbac.ClusterRoleBindingList{
			Items: []rbac.ClusterRoleBinding{
				{
					Subjects: []rbac.Subject{
						{Kind: cage_k8s.KindGroup, Name: ServiceAccountGroupAllNamespace},
					},
				},
			},
		}
		mockClientset.ClusterRoleBindings.EXPECT().List(meta.ListOptions{}).Return(&bindings, nil)

		list, err := cage_k8s_identity.ClusterRoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		require.Exactly(t, cage_k8s.KindGroup, list.Items[0].Kind)
		require.Exactly(t, ServiceAccountGroupAllNamespace, list.Items[0].Name)
	})

	t.Run("name miss", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindUser,
			Name: DoesNotExist,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		bindings := rbac.ClusterRoleBindingList{Items: []rbac.ClusterRoleBinding{}}
		mockClientset.ClusterRoleBindings.EXPECT().List(meta.ListOptions{}).Return(&bindings, nil)

		list, err := cage_k8s_identity.ClusterRoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 0)
	})

	t.Run("namespace miss", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind:      cage_k8s.KindUser,
			Name:      ServiceAccountUsername,
			Namespace: Namespace,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		var nonSut *core.Namespace
		mockClientset.Namespaces.EXPECT().Get(Namespace).Return(nonSut, Exists, nil)

		bindings := rbac.ClusterRoleBindingList{Items: []rbac.ClusterRoleBinding{}}
		mockClientset.ClusterRoleBindings.EXPECT().List(meta.ListOptions{}).Return(&bindings, nil)

		list, err := cage_k8s_identity.ClusterRoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 0)
	})

	t.Run("service account namespace unknown", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindUser,
			Name: ServiceAccountUsername,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		var nonSut *core.Namespace
		mockClientset.Namespaces.EXPECT().Get(Namespace).Return(nonSut, NotExists, nil)

		list, err := cage_k8s_identity.ClusterRoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.Nil(t, list)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("service account [%s] namespace [%s] not found", ServiceAccountUsername, Namespace))
	})

	// Assert that an error is returned if a service account's namespace does not match the query's.
	t.Run("service account namespace mismatch", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind:      cage_k8s.KindUser,
			Name:      ServiceAccountUsername,
			Namespace: DoesNotExist,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		list, err := cage_k8s_identity.ClusterRoleSubjectQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.Nil(t, list)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("query's namespace [%s] does not match query service account [%s]'s namespace [%s] ", DoesNotExist, ServiceAccountUsername, Namespace))
	})
}

func TestServiceAccountUserQuerier(t *testing.T) {
	// Assert that if the query does not require a specific object kind, still support the query
	// and filter on name/namespace.
	t.Run("compatible with empty kind", func(t *testing.T) {
		require.True(t, cage_k8s_identity.ServiceAccountUserQuerier{}.Compatible(&cage_k8s_identity.Query{}))
	})

	t.Run("compatible with subject kind", func(t *testing.T) {
		require.True(t, cage_k8s_identity.ServiceAccountUserQuerier{}.Compatible(&cage_k8s_identity.Query{Kind: cage_k8s.KindServiceAccount}))
		require.True(t, cage_k8s_identity.ServiceAccountUserQuerier{}.Compatible(&cage_k8s_identity.Query{Kind: cage_k8s.KindUser}))
		require.False(t, cage_k8s_identity.ServiceAccountUserQuerier{}.Compatible(&cage_k8s_identity.Query{Kind: cage_k8s.KindGroup}))
	})

	t.Run("user kind hit", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindUser,
			Name: ServiceAccountUsername,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		var nonSut *core.Namespace
		mockClientset.Namespaces.EXPECT().Get(Namespace).Return(nonSut, Exists, nil)

		bindings := core.ServiceAccountList{
			Items: []core.ServiceAccount{
				{
					TypeMeta:   meta.TypeMeta{Kind: cage_k8s.KindServiceAccount},
					ObjectMeta: meta.ObjectMeta{Name: ServiceAccountUsername, Namespace: Namespace},
				},
			},
		}
		mockClientset.ServiceAccounts.EXPECT().
			List(Namespace, meta.ListOptions{FieldSelector: "metadata.name=" + ServiceAccountUsernameBase}).
			Return(&bindings, nil)

		list, err := cage_k8s_identity.ServiceAccountUserQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		require.Exactly(t, cage_k8s.KindServiceAccount, list.Items[0].Kind)
		require.Exactly(t, ServiceAccountUsername, list.Items[0].Name)
	})

	t.Run("service account kind hit", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindServiceAccount,
			Name: ServiceAccountUsername,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		var nonSut *core.Namespace
		mockClientset.Namespaces.EXPECT().Get(Namespace).Return(nonSut, Exists, nil)

		bindings := core.ServiceAccountList{
			Items: []core.ServiceAccount{
				{
					TypeMeta:   meta.TypeMeta{Kind: cage_k8s.KindServiceAccount},
					ObjectMeta: meta.ObjectMeta{Name: ServiceAccountUsername, Namespace: Namespace},
				},
			},
		}
		mockClientset.ServiceAccounts.EXPECT().
			List(Namespace, meta.ListOptions{FieldSelector: "metadata.name=" + ServiceAccountUsernameBase}).
			Return(&bindings, nil)

		list, err := cage_k8s_identity.ServiceAccountUserQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		require.Exactly(t, cage_k8s.KindServiceAccount, list.Items[0].Kind)
		require.Exactly(t, ServiceAccountUsername, list.Items[0].Name)
	})

	t.Run("service account group miss", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindGroup,
			Name: ServiceAccountUsername,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		var nonSut *core.Namespace
		mockClientset.Namespaces.EXPECT().Get(Namespace).Return(nonSut, Exists, nil)

		bindings := core.ServiceAccountList{Items: []core.ServiceAccount{}}
		mockClientset.ServiceAccounts.EXPECT().
			List(Namespace, meta.ListOptions{FieldSelector: "metadata.name=" + ServiceAccountUsernameBase}).
			Return(&bindings, nil)

		list, err := cage_k8s_identity.ServiceAccountUserQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 0)
	})

	t.Run("name miss", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindServiceAccount,
			Name: ServiceAccountUsernameMiss,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		var nonSut *core.Namespace
		mockClientset.Namespaces.EXPECT().Get(Namespace).Return(nonSut, Exists, nil)

		bindings := core.ServiceAccountList{Items: []core.ServiceAccount{}}
		mockClientset.ServiceAccounts.EXPECT().
			List(Namespace, meta.ListOptions{FieldSelector: "metadata.name=" + DoesNotExist}).
			Return(&bindings, nil)

		list, err := cage_k8s_identity.ServiceAccountUserQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 0)
	})

	t.Run("namespace miss", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind:      cage_k8s.KindServiceAccount,
			Name:      ServiceAccountUsername,
			Namespace: Namespace,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		var nonSut *core.Namespace
		mockClientset.Namespaces.EXPECT().Get(Namespace).Return(nonSut, Exists, nil)

		bindings := core.ServiceAccountList{Items: []core.ServiceAccount{}}
		mockClientset.ServiceAccounts.EXPECT().
			List(Namespace, meta.ListOptions{FieldSelector: "metadata.name=" + ServiceAccountUsernameBase}).
			Return(&bindings, nil)

		list, err := cage_k8s_identity.ServiceAccountUserQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 0)
	})

	t.Run("service account namespace unknown", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindServiceAccount,
			Name: ServiceAccountUsernameNamespaceMiss,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		var nonSut *core.Namespace
		mockClientset.Namespaces.EXPECT().Get(DoesNotExist).Return(nonSut, NotExists, nil)

		list, err := cage_k8s_identity.ServiceAccountUserQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.Nil(t, list)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("service account [%s] namespace [%s] not found", ServiceAccountUsernameNamespaceMiss, DoesNotExist))
	})

	// Assert that an error is returned if a service account's namespace does not match the query's.
	t.Run("service account namespace mismatch", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind:      cage_k8s.KindServiceAccount,
			Name:      ServiceAccountUsername,
			Namespace: DoesNotExist,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		list, err := cage_k8s_identity.ServiceAccountUserQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.Nil(t, list)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("query's namespace [%s] does not match query service account [%s]'s namespace [%s] ", DoesNotExist, ServiceAccountUsername, Namespace))
	})
}

func TestServiceAccountGroupQuerier(t *testing.T) {
	// Assert that if the query does not require a specific object kind, still support the query
	// and filter on name/namespace.
	t.Run("compatible with empty kind", func(t *testing.T) {
		require.True(t, cage_k8s_identity.ServiceAccountGroupQuerier{}.Compatible(&cage_k8s_identity.Query{}))
	})

	t.Run("compatible with subject kind", func(t *testing.T) {
		require.False(t, cage_k8s_identity.ServiceAccountGroupQuerier{}.Compatible(&cage_k8s_identity.Query{Kind: cage_k8s.KindServiceAccount}))
		require.False(t, cage_k8s_identity.ServiceAccountGroupQuerier{}.Compatible(&cage_k8s_identity.Query{Kind: cage_k8s.KindUser}))
		require.True(t, cage_k8s_identity.ServiceAccountGroupQuerier{}.Compatible(&cage_k8s_identity.Query{Kind: cage_k8s.KindGroup}))
	})

	t.Run("hit", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind: cage_k8s.KindGroup,
			Name: ServiceAccountGroupOneNamespace,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		var nonSut *core.Namespace
		mockClientset.Namespaces.EXPECT().Get(Namespace).Return(nonSut, Exists, nil)

		list, err := cage_k8s_identity.ServiceAccountGroupQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		require.Exactly(t, cage_k8s.KindGroup, list.Items[0].Kind)
		require.Exactly(t, ServiceAccountGroupOneNamespace, list.Items[0].Name)
	})

	t.Run("namespace unknown", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind:      cage_k8s.KindGroup,
			Name:      ServiceAccountGroupNamespaceMiss,
			Namespace: DoesNotExist,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		var nonSut *core.Namespace
		mockClientset.Namespaces.EXPECT().Get(DoesNotExist).Return(nonSut, NotExists, nil)

		list, err := cage_k8s_identity.ServiceAccountGroupQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.Nil(t, list)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("service account [%s] namespace [%s] not found", ServiceAccountGroupNamespaceMiss, DoesNotExist))
	})

	// Assert that an error is returned if a service account's namespace does not match the query's.
	t.Run("namespace mismatch", func(t *testing.T) {
		query := cage_k8s_identity.Query{
			Kind:      cage_k8s.KindGroup,
			Name:      ServiceAccountGroupOneNamespace,
			Namespace: DoesNotExist,
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockClientset := mock_core.NewClientset(mockCtrl)

		list, err := cage_k8s_identity.ServiceAccountGroupQuerier{}.Do(ctx(), mockClientset.ToReal(), &query)
		require.Nil(t, list)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("query's namespace [%s] does not match query service account [%s]'s namespace [%s] ", DoesNotExist, ServiceAccountGroupOneNamespace, Namespace))
	})
}
