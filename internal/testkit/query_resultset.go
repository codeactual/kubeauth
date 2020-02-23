package testkit

import (
	cage_k8s_identity "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/identity"
)

type QueryResultset struct {
	CoreGroup           *cage_k8s_identity.IdentityList
	CoreUser            *cage_k8s_identity.IdentityList
	RoleSubject         *cage_k8s_identity.IdentityList
	ClusterRoleSubject  *cage_k8s_identity.IdentityList
	ServiceAccountUser  *cage_k8s_identity.IdentityList
	ServiceAccountGroup *cage_k8s_identity.IdentityList
	ConfigUser          *cage_k8s_identity.IdentityList
}

func NewQueryResultset() QueryResultset {
	return QueryResultset{
		CoreGroup:           &cage_k8s_identity.IdentityList{},
		CoreUser:            &cage_k8s_identity.IdentityList{},
		RoleSubject:         &cage_k8s_identity.IdentityList{},
		ClusterRoleSubject:  &cage_k8s_identity.IdentityList{},
		ServiceAccountUser:  &cage_k8s_identity.IdentityList{},
		ServiceAccountGroup: &cage_k8s_identity.IdentityList{},
		ConfigUser:          &cage_k8s_identity.IdentityList{},
	}
}
