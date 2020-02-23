// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

//go:generate mockgen -copyright_file=$LICENSE_HEADER -package=mock -destination=$GODIR/mock/querier.go -source=$GODIR/$GOFILE
package identity

import (
	"context"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pkg/errors"

	cage_k8s "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1"
	cage_k8s_core "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/core"
	cage_k8s_rbac "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac"
)

// Querier implementations perform queries of an identity-related object kind, e.g. ServiceAccount.
//
// This decomposition makes a trade-off between more Go types/files and the ability to define the
// sub-queries independently in a more maintainable way.
type Querier interface {
	// String returns a unique description of the type of result provided by the querier.
	String() string

	// Compatible returns true if the implementation can serve the query.
	//
	// For example, a query may specify an object kind, e.g. ClusterRole, but a querier may not
	// know how to query it because it only supports Group.
	Compatible(*Query) bool

	// Do performs the query.
	Do(context.Context, *cage_k8s_core.Clientset, *Query) (*IdentityList, error)
}

// ConfigUserQuerier queries a kubectl context for its user value.
type ConfigUserQuerier struct{}

// String returns a unique description of the type of result provided by the querier.
//
// It implements Querier.
func (q ConfigUserQuerier) String() string {
	return "kubeconfig context"
}

// Compatible returns true if the implementation can serve the query.
//
// It implements Querier.
func (q ConfigUserQuerier) Compatible(query *Query) bool {
	return (query.Kind == "" || query.Kind == cage_k8s.KindUser)
}

// Do performs the query.
//
// It implements Querier.
func (q ConfigUserQuerier) Do(ctx context.Context, _ *cage_k8s_core.Clientset, query *Query) (*IdentityList, error) {
	if query.ClientCmdConfig == nil {
		return nil, errors.Errorf("[%s] querier received a nil kubeconfig object", q)
	}

	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "query cancelled")
	default:
	}

	var found IdentityList

	for _, contextObj := range query.ClientCmdConfig.Contexts {
		if query.Namespace != "" && query.Namespace != contextObj.Namespace {
			continue
		}

		if query.Name == "" || query.Name == contextObj.AuthInfo {
			found.Items = append(found.Items, Identity{
				TypeMeta:   meta.TypeMeta{Kind: cage_k8s.KindUser},
				ObjectMeta: meta.ObjectMeta{Name: contextObj.AuthInfo, Namespace: contextObj.Namespace},
			})
		}
	}

	return &found, nil
}

var _ Querier = (*ConfigUserQuerier)(nil)

// CoreGroupQuerier queries a hard-coded set of group names enumerated in the API server source code.
//
// They're included as string literals instead of imported constants in order to avoid k8s.io/apiserver
// and its transitive dependencies.
//
// https://github.com/kubernetes/apiserver/blob/kubernetes-1.17.0/pkg/authentication/user/user.go#L69
type CoreGroupQuerier struct{}

// String returns a unique description of the type of result provided by the querier.
//
// It implements Querier.
func (q CoreGroupQuerier) String() string {
	return "system-defined group"
}

// Compatible returns true if the implementation can serve the query.
//
// It implements Querier.
func (q CoreGroupQuerier) Compatible(query *Query) bool {
	return query.Kind == "" || query.Kind == cage_k8s.KindGroup
}

// Do performs the query.
//
// It implements Querier.
func (q CoreGroupQuerier) Do(ctx context.Context, _ *cage_k8s_core.Clientset, query *Query) (*IdentityList, error) {
	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "query cancelled")
	default:
	}

	ids := []Identity{
		{TypeMeta: meta.TypeMeta{Kind: cage_k8s.KindGroup}, ObjectMeta: meta.ObjectMeta{Name: "system:masters"}},
		{TypeMeta: meta.TypeMeta{Kind: cage_k8s.KindGroup}, ObjectMeta: meta.ObjectMeta{Name: "system:nodes"}},
		{TypeMeta: meta.TypeMeta{Kind: cage_k8s.KindGroup}, ObjectMeta: meta.ObjectMeta{Name: "system:unauthenticated"}},
		{TypeMeta: meta.TypeMeta{Kind: cage_k8s.KindGroup}, ObjectMeta: meta.ObjectMeta{Name: "system:authenticated"}},
	}

	var found IdentityList

	for _, id := range ids {
		if query.Name == "" || query.Name == id.Name {
			found.Items = append(found.Items, id)
		}
	}

	return &found, nil
}

var _ Querier = (*CoreGroupQuerier)(nil)

// CoreGroupQuerier queries a hard-coded set of user names enumerated in the API server source code.
//
// They're included as string literals instead of imported constants in order to avoid k8s.io/apiserver
// and its transitive dependencies.
//
// https://github.com/kubernetes/apiserver/blob/kubernetes-1.17.0/pkg/authentication/user/user.go#L69
type CoreUserQuerier struct{}

// String returns a unique description of the type of result provided by the querier.
//
// It implements Querier.
func (q CoreUserQuerier) String() string {
	return "system-defined user"
}

// Compatible returns true if the implementation can serve the query.
//
// It implements Querier.
func (q CoreUserQuerier) Compatible(query *Query) bool {
	return query.Kind == "" || query.Kind == cage_k8s.KindUser
}

// Do performs the query.
//
// It implements Querier.
func (q CoreUserQuerier) Do(ctx context.Context, _ *cage_k8s_core.Clientset, query *Query) (*IdentityList, error) {
	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "query cancelled")
	default:
	}

	ids := []Identity{
		{TypeMeta: meta.TypeMeta{Kind: cage_k8s.KindUser}, ObjectMeta: meta.ObjectMeta{Name: "system:anonymous"}},
		{TypeMeta: meta.TypeMeta{Kind: cage_k8s.KindUser}, ObjectMeta: meta.ObjectMeta{Name: "system:apiserver"}},
		{TypeMeta: meta.TypeMeta{Kind: cage_k8s.KindUser}, ObjectMeta: meta.ObjectMeta{Name: "system:kube-proxy"}},
		{TypeMeta: meta.TypeMeta{Kind: cage_k8s.KindUser}, ObjectMeta: meta.ObjectMeta{Name: "system:kube-controller-manager"}},
		{TypeMeta: meta.TypeMeta{Kind: cage_k8s.KindUser}, ObjectMeta: meta.ObjectMeta{Name: "system:kube-scheduler"}},
	}

	var found IdentityList

	for _, id := range ids {
		if query.Name == "" || query.Name == id.Name {
			found.Items = append(found.Items, id)
		}
	}

	return &found, nil
}

var _ Querier = (*CoreUserQuerier)(nil)

// RoleSubjectQuerier queries the API for role subjects.
type RoleSubjectQuerier struct{}

// String returns a unique description of the type of result provided by the querier.
//
// It implements Querier.
func (q RoleSubjectQuerier) String() string {
	return "role binding subject"
}

// Compatible returns true if the implementation can serve the query.
//
// It implements Querier.
func (q RoleSubjectQuerier) Compatible(query *Query) bool {
	return query.Kind == "" || query.Kind == cage_k8s.KindUser || query.Kind == cage_k8s.KindGroup
}

// Do performs the query.
//
// It implements Querier.
func (q RoleSubjectQuerier) Do(ctx context.Context, clientset *cage_k8s_core.Clientset, query *Query) (*IdentityList, error) {
	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "query cancelled")
	default:
	}

	// Detect names of service accounts and parse them when found. If found and the name includes a namespace,
	// verify the namespace exists.

	querySaNamespace, querySaName, querySaIsGroup, querySaIsValid := cage_k8s_rbac.ParseServiceAccount(query.Name)
	if querySaIsValid {
		if querySaNamespace != "" {
			if query.Namespace != "" && querySaNamespace != query.Namespace {
				return nil, errors.Errorf("query's namespace [%s] does not match query service account [%s]'s namespace [%s] ", query.Namespace, query.Name, querySaNamespace)
			}

			_, exists, err := clientset.Namespaces.Get(querySaNamespace)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if !exists {
				return nil, errors.Errorf("service account [%s] namespace [%s] not found", query.Name, querySaNamespace)
			}
		}
	}

	// Scan role bindings for subjects which match the queried name. Apply queried namespace if provided.

	res, err := clientset.RoleBindings.List(query.Namespace, meta.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get role binding list")
	}

	var list IdentityList
	for _, r := range res.Items {
		for _, s := range r.Subjects {
			if query.Namespace != "" && query.Namespace != s.Namespace {
				continue
			}

			var match bool

			if query.Name == "" {
				match = true
			} else {
				if querySaIsValid {
					if querySaIsGroup {
						match = s.Kind == cage_k8s.KindGroup && s.Name == query.Name
					} else {
						match = s.Kind == cage_k8s.KindServiceAccount && s.Name == querySaName
					}
				} else {
					match = (query.Kind == "" || s.Kind == query.Kind) && s.Name == query.Name
				}
			}

			if match {
				list.Items = append(list.Items, Identity{
					ObjectMeta: meta.ObjectMeta{
						Name:      s.Name,
						Namespace: s.Namespace,
					},
					TypeMeta: meta.TypeMeta{
						Kind: s.Kind,
					},
					Source: &IdentitySource{
						ObjectMeta: meta.ObjectMeta{
							Name:      r.Name,
							Namespace: r.Namespace,
						},
						TypeMeta: meta.TypeMeta{
							Kind: cage_k8s.KindRoleBinding, // as of v0.16.4, r.Kind is empty
						},
					},
				})
			}
		}
	}

	return &list, nil
}

var _ Querier = (*RoleSubjectQuerier)(nil)

// ClusterRoleSubjectQuerier queries the API for role subjects.
type ClusterRoleSubjectQuerier struct{}

// String returns a unique description of the type of result provided by the querier.
//
// It implements Querier.
func (q ClusterRoleSubjectQuerier) String() string {
	return "cluster role binding subject"
}

// Compatible returns true if the implementation can serve the query.
//
// It implements Querier.
func (q ClusterRoleSubjectQuerier) Compatible(query *Query) bool {
	return query.Kind == "" || query.Kind == cage_k8s.KindUser || query.Kind == cage_k8s.KindGroup
}

// Do performs the query.
//
// It implements Querier.
func (q ClusterRoleSubjectQuerier) Do(ctx context.Context, clientset *cage_k8s_core.Clientset, query *Query) (*IdentityList, error) {
	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "query cancelled")
	default:
	}

	// Detect names of service accounts and parse them when found. If found and the name includes a namespace,
	// verify the namespace exists.

	querySaNamespace, querySaName, querySaIsGroup, querySaIsValid := cage_k8s_rbac.ParseServiceAccount(query.Name)
	if querySaIsValid {
		if querySaNamespace != "" {
			if query.Namespace != "" && querySaNamespace != query.Namespace {
				return nil, errors.Errorf("query's namespace [%s] does not match query service account [%s]'s namespace [%s] ", query.Namespace, query.Name, querySaNamespace)
			}

			_, exists, err := clientset.Namespaces.Get(querySaNamespace)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if !exists {
				return nil, errors.Errorf("service account [%s] namespace [%s] not found", query.Name, querySaNamespace)
			}
		}
	}

	// Scan role bindings for subjects which match the queried name. Apply queried namespace if provided.

	res, err := clientset.ClusterRoleBindings.List(meta.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cluster role binding list")
	}

	var list IdentityList
	for _, r := range res.Items {
		for _, s := range r.Subjects {
			if query.Namespace != "" && query.Namespace != s.Namespace {
				continue
			}

			var match bool

			if query.Name == "" {
				match = true
			} else {
				if querySaIsValid {
					if querySaIsGroup {
						match = s.Kind == cage_k8s.KindGroup && s.Name == query.Name
					} else {
						match = s.Kind == cage_k8s.KindServiceAccount && s.Name == querySaName
					}
				} else {
					match = (query.Kind == "" || s.Kind == query.Kind) && s.Name == query.Name
				}
			}

			if match {
				list.Items = append(list.Items, Identity{
					ObjectMeta: meta.ObjectMeta{
						Name:      s.Name,
						Namespace: s.Namespace,
					},
					TypeMeta: meta.TypeMeta{
						Kind: s.Kind,
					},
					Source: &IdentitySource{
						ObjectMeta: meta.ObjectMeta{
							Name:      r.Name,
							Namespace: r.Namespace,
						},
						TypeMeta: meta.TypeMeta{
							Kind: cage_k8s.KindClusterRoleBinding, // as of v0.16.4, r.Kind is empty
						},
					},
				})
			}
		}
	}

	return &list, nil
}

var _ Querier = (*ClusterRoleSubjectQuerier)(nil)

// ServiceAccountUserQuerier queries the API for service account based users.
type ServiceAccountUserQuerier struct{}

// String returns a unique description of the type of result provided by the querier.
//
// It implements Querier.
func (q ServiceAccountUserQuerier) String() string {
	return "service account based user"
}

// Compatible returns true if the implementation can serve the query.
//
// It implements Querier.
func (q ServiceAccountUserQuerier) Compatible(query *Query) bool {
	return query.Kind == "" || query.Kind == cage_k8s.KindUser || query.Kind == cage_k8s.KindServiceAccount
}

// Do performs the query.
//
// It implements Querier.
func (q ServiceAccountUserQuerier) Do(ctx context.Context, clientset *cage_k8s_core.Clientset, query *Query) (*IdentityList, error) {
	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "query cancelled")
	default:
	}

	var list IdentityList
	var listOpts meta.ListOptions

	querySaNamespace, querySaName, querySaIsGroup, querySaIsValid := cage_k8s_rbac.ParseServiceAccount(query.Name)
	if !querySaIsValid || querySaIsGroup {
		return &list, nil
	}

	if querySaNamespace != "" {
		if query.Namespace != "" && querySaNamespace != query.Namespace {
			return nil, errors.Errorf("query's namespace [%s] does not match query service account [%s]'s namespace [%s] ", query.Namespace, query.Name, querySaNamespace)
		}

		_, exists, err := clientset.Namespaces.Get(querySaNamespace)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if !exists {
			return nil, errors.Errorf("service account [%s] namespace [%s] not found", query.Name, querySaNamespace)
		}
	}

	listOpts.FieldSelector = "metadata.name=" + querySaName

	res, err := clientset.ServiceAccounts.List(querySaNamespace, listOpts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get role list")
	}

	for _, s := range res.Items {
		list.Items = append(list.Items, Identity{
			ObjectMeta: meta.ObjectMeta{
				Name:      s.Name,
				Namespace: s.Namespace,
			},
			TypeMeta: meta.TypeMeta{
				Kind: cage_k8s.KindServiceAccount,
			},
		})
	}

	return &list, nil
}

var _ Querier = (*ServiceAccountUserQuerier)(nil)

// ServiceAccountGroupQuerier detects valid names of service account based groups
// and queries the API to validate their namespaces if needed. If all validation
// checks pass, the group is returned in the identity list.
type ServiceAccountGroupQuerier struct{}

// String returns a unique description of the type of result provided by the querier.
//
// It implements Querier.
func (q ServiceAccountGroupQuerier) String() string {
	return "service account based group"
}

// Compatible returns true if the implementation can serve the query.
//
// It implements Querier.
func (q ServiceAccountGroupQuerier) Compatible(query *Query) bool {
	return query.Kind == "" || query.Kind == cage_k8s.KindGroup
}

// Do performs the query.
//
// It implements Querier.
func (q ServiceAccountGroupQuerier) Do(ctx context.Context, clientset *cage_k8s_core.Clientset, query *Query) (*IdentityList, error) {
	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "query cancelled")
	default:
	}

	var list IdentityList

	if query.Name == "" {
		return &list, nil
	}

	querySaNamespace, _, querySaIsGroup, querySaIsValid := cage_k8s_rbac.ParseServiceAccount(query.Name)
	if !querySaIsValid || !querySaIsGroup {
		return &list, nil
	}

	if querySaNamespace != "" {
		if query.Namespace != "" && querySaNamespace != query.Namespace {
			return nil, errors.Errorf("query's namespace [%s] does not match query service account [%s]'s namespace [%s] ", query.Namespace, query.Name, querySaNamespace)
		}

		_, exists, err := clientset.Namespaces.Get(querySaNamespace)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if !exists {
			return nil, errors.Errorf("service account [%s] namespace [%s] not found", query.Name, querySaNamespace)
		}
	}

	list.Items = append(list.Items, Identity{
		ObjectMeta: meta.ObjectMeta{
			Name:      query.Name,
			Namespace: querySaNamespace,
		},
		TypeMeta: meta.TypeMeta{
			Kind: cage_k8s.KindGroup,
		},
	})

	return &list, nil
}

var _ Querier = (*ServiceAccountGroupQuerier)(nil)
