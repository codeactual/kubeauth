// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package identity

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	cage_k8s_core "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/core"
)

type Registry struct {
	CoreGroup           Querier
	CoreUser            Querier
	RoleSubject         Querier
	ClusterRoleSubject  Querier
	ServiceAccountUser  Querier
	ServiceAccountGroup Querier
	ConfigUser          Querier

	Clientset *cage_k8s_core.Clientset
}

// NewRegistry builds a registry of known and discovered users.
func NewRegistry(clientset *cage_k8s_core.Clientset) *Registry {
	return &Registry{
		CoreGroup:           CoreGroupQuerier{},
		CoreUser:            CoreUserQuerier{},
		RoleSubject:         RoleSubjectQuerier{},
		ClusterRoleSubject:  ClusterRoleSubjectQuerier{},
		ServiceAccountUser:  ServiceAccountUserQuerier{},
		ServiceAccountGroup: ServiceAccountGroupQuerier{},
		ConfigUser:          ConfigUserQuerier{},
		Clientset:           clientset,
	}
}

func (reg *Registry) Query(ctx context.Context, options ...QueryOption) (*IdentityList, error) {
	query := NewQuery(options...)

	// Initialize all available queriers.

	queriers := []Querier{
		reg.CoreGroup,
		reg.CoreUser,
		reg.RoleSubject,
		reg.ClusterRoleSubject,
		reg.ServiceAccountUser,
		reg.ServiceAccountGroup,
	}

	if query.ClientCmdConfig != nil {
		queriers = append(queriers, reg.ConfigUser)
	}

	// Run compatible queriers in parallel.

	g, gCtx := errgroup.WithContext(ctx)
	var fullList IdentityList
	var mu sync.Mutex

	// Use a constructor instead of a function literal in errgroup.Group.Go calls
	// to avoid accidental closure value issues.
	newErrGroupFn := func(querier Querier) func() error {
		return func() error {
			querierType := querier.String()

			list, err := querier.Do(gCtx, reg.Clientset, query)
			if err != nil {
				return errors.Wrapf(err, "identity querier [%s] did not finish", querierType)
			}

			for _, item := range list.Items {
				item.Querier = querierType
				mu.Lock()
				fullList.Items = append(fullList.Items, item)
				mu.Unlock()
			}

			return nil
		}
	}

	for _, querier := range queriers {
		if !querier.Compatible(query) {
			continue
		}
		g.Go(newErrGroupFn(querier))
	}

	if err := g.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	return &fullList, nil
}
