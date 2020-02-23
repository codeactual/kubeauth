// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rbac

import (
	"strings"

	"github.com/pkg/errors"
)

type BindingSelector struct {
	RoleName    string
	BindingName string
}

func NewBindingSelector(s string) (*BindingSelector, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, errors.Errorf("selector [%s] does not use format <role name>:<binding name>", s)
	}
	return &BindingSelector{RoleName: parts[0], BindingName: parts[1]}, nil
}

// ParseServiceAccountUser parses a service account user name.
//
// For system:serviceaccount:a:b, it returns namespace "a" and basename "b".
func ParseServiceAccountUser(user string) (namespace, basename string, _ error) {
	parts := strings.Split(user, ":")

	if len(parts) != 4 {
		return "", "", errors.Errorf("service account user [%s] contains [%d] parts, expected 4", user, len(parts))
	}

	if parts[0] != "system" || parts[1] != "serviceaccount" {
		return "", "", errors.Errorf("service account user [%s] does not begin with 'system:serviceaccount:'", user)
	}

	return parts[2], parts[3], nil
}

// ParseServiceAccountGroup parses a service account group name.
//
// For system:serviceaccounts, it returns namespace "".
// For system:serviceaccounts:a, it returns namespace "a".
func ParseServiceAccountGroup(group string) (namespace string, _ error) {
	parts := strings.Split(group, ":")

	if len(parts) != 2 && len(parts) != 3 {
		return "", errors.Errorf("service account group [%s] contains [%d] parts, expected 2 or 3", group, len(parts))
	}

	if parts[0] != "system" || parts[1] != "serviceaccounts" {
		return "", errors.Errorf("service account group [%s] does not begin with 'system:serviceaccounts:'", group)
	}

	if len(parts) == 2 {
		return "", nil
	}

	return parts[2], nil
}

// ParseServiceAccount parses a service account name and identifies whether its a user or group.
//
// For user name system:serviceaccount:a:b, it returns basename "a" and namespace "b".
// For group name system:serviceaccounts:a, it returns basename "" and namespace "a".
func ParseServiceAccount(n string) (namespace, basename string, isGroup, isValid bool) {
	namespace, basename, _ = ParseServiceAccountUser(n)
	if basename == "" {
		namespace, err := ParseServiceAccountGroup(n)
		if err == nil {
			return namespace, "", true, true
		}
	} else {
		return namespace, basename, false, true
	}

	return "", "", false, false
}
