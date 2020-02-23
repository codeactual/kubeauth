// Copyright (C) 2020 The kubeauth Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gomock

// Reminder: functions which build matchers must return the gomock.Matcher type in order for
// gomock (as of post-1.3.1 817c00c) to recognize the value as a matcher. If an implementation
// type is returned, gomock will simply compare equality.

import (
	"context"
	"fmt"
	"regexp"

	"github.com/golang/mock/gomock"
)

// NonSut may be clearer than nil as a Return parameter in order to "document"
// that parameter as outside the SUT in the test case.
func NonSut() interface{} {
	return nil
}

type matchErrShortRegexp struct {
	re *regexp.Regexp
}

func (m *matchErrShortRegexp) Matches(x interface{}) bool {
	if err, ok := x.(error); ok {
		if err == nil {
			return false
		}

		return m.re.MatchString(err.Error())
	}

	return false
}

func (m *matchErrShortRegexp) String() string {
	return fmt.Sprintf("error's Error() string matches regexp `%s`", m.re)
}

var _ gomock.Matcher = (*matchErrShortRegexp)(nil)

// ErrShortRegexp compares the regular expression against an expected error's Error() string.
func ErrShortRegexp(re *regexp.Regexp) gomock.Matcher { // must return this type for gomock to recognize it
	return &matchErrShortRegexp{re: re}
}

type matchContextNonNil struct {
}

func (m *matchContextNonNil) Matches(x interface{}) bool {
	if ctx, ok := x.(context.Context); ok && ctx != nil {
		return true
	}

	return false
}

func (m *matchContextNonNil) String() string {
	return "context implementation is non nil"
}

var _ gomock.Matcher = (*matchContextNonNil)(nil)

// ContextNonNil checks if the value is a context.Context and is non-nil.
func ContextNonNil() gomock.Matcher { // must return this type for gomock to recognize it
	return &matchContextNonNil{}
}
