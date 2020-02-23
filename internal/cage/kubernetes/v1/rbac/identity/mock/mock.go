// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package mock

import (
	"fmt"
	"os"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/sanity-io/litter"

	cage_k8s "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1"
	cage_k8s_core "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/core"
	cage_k8s_identity "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/identity"
)

type Registry struct {
	CoreGroup           *MockQuerier
	CoreUser            *MockQuerier
	RoleSubject         *MockQuerier
	ClusterRoleSubject  *MockQuerier
	ServiceAccountUser  *MockQuerier
	ServiceAccountGroup *MockQuerier
	ConfigUser          *MockQuerier
}

func NewRegistry(mockCtrl *gomock.Controller) *Registry {
	r := Registry{
		CoreGroup:           NewMockQuerier(mockCtrl),
		CoreUser:            NewMockQuerier(mockCtrl),
		RoleSubject:         NewMockQuerier(mockCtrl),
		ClusterRoleSubject:  NewMockQuerier(mockCtrl),
		ServiceAccountUser:  NewMockQuerier(mockCtrl),
		ServiceAccountGroup: NewMockQuerier(mockCtrl),
		ConfigUser:          NewMockQuerier(mockCtrl),
	}

	// Reuse real Querier.Compatible checks because they're currently fast.
	r.CoreGroup.EXPECT().Compatible(gomock.Any()).DoAndReturn(func(query interface{}) bool {
		return cage_k8s_identity.CoreGroupQuerier{}.Compatible(query.(*cage_k8s_identity.Query))
	}).AnyTimes()
	r.CoreUser.EXPECT().Compatible(gomock.Any()).DoAndReturn(func(query interface{}) bool {
		return cage_k8s_identity.CoreUserQuerier{}.Compatible(query.(*cage_k8s_identity.Query))
	}).AnyTimes()
	r.RoleSubject.EXPECT().Compatible(gomock.Any()).DoAndReturn(func(query interface{}) bool {
		return cage_k8s_identity.RoleSubjectQuerier{}.Compatible(query.(*cage_k8s_identity.Query))
	}).AnyTimes()
	r.ClusterRoleSubject.EXPECT().Compatible(gomock.Any()).DoAndReturn(func(query interface{}) bool {
		return cage_k8s_identity.ClusterRoleSubjectQuerier{}.Compatible(query.(*cage_k8s_identity.Query))
	}).AnyTimes()
	r.ServiceAccountUser.EXPECT().Compatible(gomock.Any()).DoAndReturn(func(query interface{}) bool {
		return cage_k8s_identity.ServiceAccountUserQuerier{}.Compatible(query.(*cage_k8s_identity.Query))
	}).AnyTimes()
	r.ServiceAccountGroup.EXPECT().Compatible(gomock.Any()).DoAndReturn(func(query interface{}) bool {
		return cage_k8s_identity.ServiceAccountGroupQuerier{}.Compatible(query.(*cage_k8s_identity.Query))
	}).AnyTimes()
	r.ConfigUser.EXPECT().Compatible(gomock.Any()).DoAndReturn(func(query interface{}) bool {
		return cage_k8s_identity.ConfigUserQuerier{}.Compatible(query.(*cage_k8s_identity.Query))
	}).AnyTimes()

	// Reuse real Querier.String values so they're available to assert which querier
	// provided an Identity to the Registry, and also to verbose messages emitted by
	// the commands during "go test -v".
	r.CoreGroup.EXPECT().String().DoAndReturn(func() string {
		return cage_k8s_identity.CoreGroupQuerier{}.String()
	}).AnyTimes()
	r.CoreUser.EXPECT().String().DoAndReturn(func() string {
		return cage_k8s_identity.CoreUserQuerier{}.String()
	}).AnyTimes()
	r.RoleSubject.EXPECT().String().DoAndReturn(func() string {
		return cage_k8s_identity.RoleSubjectQuerier{}.String()
	}).AnyTimes()
	r.ClusterRoleSubject.EXPECT().String().DoAndReturn(func() string {
		return cage_k8s_identity.ClusterRoleSubjectQuerier{}.String()
	}).AnyTimes()
	r.ServiceAccountUser.EXPECT().String().DoAndReturn(func() string {
		return cage_k8s_identity.ServiceAccountUserQuerier{}.String()
	}).AnyTimes()
	r.ServiceAccountGroup.EXPECT().String().DoAndReturn(func() string {
		return cage_k8s_identity.ServiceAccountGroupQuerier{}.String()
	}).AnyTimes()
	r.ConfigUser.EXPECT().String().DoAndReturn(func() string {
		return cage_k8s_identity.ConfigUserQuerier{}.String()
	}).AnyTimes()

	return &r
}

func (r *Registry) ToReal(clientset *cage_k8s_core.Clientset) *cage_k8s_identity.Registry {
	return &cage_k8s_identity.Registry{
		CoreGroup:           r.CoreGroup,
		CoreUser:            r.CoreUser,
		RoleSubject:         r.RoleSubject,
		ClusterRoleSubject:  r.ClusterRoleSubject,
		ServiceAccountUser:  r.ServiceAccountUser,
		ServiceAccountGroup: r.ServiceAccountGroup,
		ConfigUser:          r.ConfigUser,
		Clientset:           clientset,
	}
}

type matchQuery struct {
	expected      *cage_k8s_identity.Query
	allNamespaces bool
}

func (m *matchQuery) Matches(x interface{}) bool {
	actual, ok := x.(*cage_k8s_identity.Query)
	matches := ok && actual != nil &&
		m.expected.Kind == actual.Kind &&
		m.expected.Name == actual.Name

	if !m.allNamespaces {
		matches = matches && m.expected.Namespace == actual.Namespace
	}

	// Only expect the kubectl config object to be present during user queries.
	if m.expected.Kind == cage_k8s.KindUser {
		matches = matches &&
			m.expected.ClientCmdConfig.CurrentContext == actual.ClientCmdConfig.CurrentContext &&
			gomock.Eq(m.expected.ClientCmdConfig.Contexts).Matches(actual.ClientCmdConfig.Contexts) &&
			gomock.Eq(m.expected.ClientCmdConfig.Clusters).Matches(actual.ClientCmdConfig.Clusters)
	}

	if !matches && testing.Verbose() {
		fmt.Fprintf(os.Stderr,
			"\nQuery mismatch:\n\n"+
				"kind: expected [%s] actual [%s]\n"+
				"namespace (--all-namespaces: %t): expected [%s] actual [%s]\n"+
				"name: expected [%s] actual [%s]\n"+
				"current context: expected [%s] actual [%s]\n"+
				"\ncontexts:\n\n"+
				"expected: %s\n"+
				"actual: %s\n"+
				"\nclusters:\n\n"+
				"expected: %s\n"+
				"actual: %s\n",
			m.expected.Kind,
			actual.Kind,
			m.allNamespaces,
			m.expected.Namespace,
			actual.Namespace,
			m.expected.Name,
			actual.Name,
			m.expected.ClientCmdConfig.CurrentContext,
			actual.ClientCmdConfig.CurrentContext,
			litter.Sdump(m.expected.ClientCmdConfig.Contexts),
			litter.Sdump(actual.ClientCmdConfig.Contexts),
			litter.Sdump(m.expected.ClientCmdConfig.Clusters),
			litter.Sdump(actual.ClientCmdConfig.Clusters),
		)
	}

	return matches
}

func (m *matchQuery) String() string {
	return "is equal to expected Query"
}

var _ gomock.Matcher = (*matchQuery)(nil)

// MatchQuery returns a matcher to assert the equality of Query fields.
func MatchQuery(allNamespaces bool, expected *cage_k8s_identity.Query) gomock.Matcher { // must return this type for gomock to recognize it
	return &matchQuery{allNamespaces: allNamespaces, expected: expected}
}
