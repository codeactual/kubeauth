// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package identity

import (
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// QueryOption implementations accept the current Query state and update it based
// on option-specific logic.
//
// It supports a functional option API based on https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis.
type QueryOption func(*Query)

// QueryNamespace limits the query scope of an RBAC related query to this namespace.
//
// To align with kubectl's --namespace/--all-namespaces behavior, if this option is set to a
// non-empty value, Querier implementations will ignore it if the queried dataset
// is namespace agnostic, e.g. cluster roles. In other words, if the selected namespace is "frontend"
// and cluster roles are queried, any cluster role will be included in query results as
// long as it matchesj other criteria.
func QueryNamespace(val string) QueryOption {
	return func(q *Query) {
		q.Namespace = val
	}
}

// QueryKind limits the query scope of an RBAC related query to a specific object kind.
func QueryKind(val string) QueryOption {
	return func(q *Query) {
		q.Kind = val
	}
}

// QueryName limits the query scope of an RBAC related query to a specific name.
func QueryName(val string) QueryOption {
	return func(q *Query) {
		q.Name = val
	}
}

// QueryClientCmdConfig expands the query scope of an RBAC related query to seek matches
// from the config's identity-related entities.
func QueryClientCmdConfig(val *clientcmdapi.Config) QueryOption {
	return func(q *Query) {
		q.ClientCmdConfig = val.DeepCopy()
	}
}

// NewQuery returns a Query initialized with all input options.
func NewQuery(options ...QueryOption) *Query {
	q := Query{}

	for _, o := range options {
		o(&q)
	}

	return &q
}

// Query holds facets which limit an RBAC related query's result set.
type Query struct {
	// Kind determines which Querier implementations are used by only running those which support this kind.
	Kind string

	// Name limits which identities are returned from Querier implementations. If it matches a candidate's
	// name, or if it is empty, the candidate is included the returned List.
	Name string

	// Namespace limits which are returned from Querier implementations. For example, if the querier
	// consumes a RoleBinding list, only bindings from the selected namespace are considered.
	Namespace string

	// ClientCmdConfig provides kubectl config values from which to seek query matches.
	ClientCmdConfig *clientcmdapi.Config
}
