// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package config_test

import (
	"bytes"
	"context"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/kubectl/config"
	cage_exec "github.com/codeactual/kubeauth/internal/cage/os/exec"
	mock_exec "github.com/codeactual/kubeauth/internal/cage/os/exec/mock"
	testkit_file "github.com/codeactual/kubeauth/internal/cage/testkit/os/file"
)

type ConfigSuite struct {
	suite.Suite
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

func (s *ConfigSuite) TestFileGetCurrentContext() {
	t := s.T()

	filename := filepath.Join(testkit_file.FixtureDataDir(), "kubeconfig-orig.yml")
	client := config.NewDefaultClient()

	file, err := client.Parse(filename)
	require.NoError(t, err)

	name, obj, err := file.GetCurrentContext()
	require.NoError(t, err)
	require.Exactly(t, "some-context", name)
	require.Exactly(t, "some-user", obj.AuthInfo)
	require.Exactly(t, "some-namespace", obj.Namespace)
}

func (s *ConfigSuite) TestFileGetCurrentCluster() {
	t := s.T()

	filename := filepath.Join(testkit_file.FixtureDataDir(), "kubeconfig-orig.yml")
	client := config.NewDefaultClient()

	file, err := client.Parse(filename)
	require.NoError(t, err)

	name, obj, err := file.GetCurrentCluster()
	require.NoError(t, err)
	require.Exactly(t, "some-cluster", name)
	require.Exactly(t, "https://1.2.3.4", obj.Server)
}

// TestClientParse samples File contents not covered by other cases.
func (s *ConfigSuite) TestClientParse() {
	t := s.T()

	filename := filepath.Join(testkit_file.FixtureDataDir(), "kubeconfig-orig.yml")
	client := config.NewDefaultClient()

	file, err := client.Parse(filename)
	require.NoError(t, err)

	require.Exactly(t, filename, file.Name)
	require.Exactly(t, "https://1.2.3.4", file.RestConfig.Host)
}

func (s *ConfigSuite) TestClientUpsertToken() {
	t := s.T()
	ctx := context.Background()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	filename := filepath.Join(testkit_file.FixtureDataDir(), "kubeconfig-orig.yml")
	client := config.NewDefaultClient()

	file, err := client.Parse(filename)
	require.NoError(t, err)

	expectToken := []byte("some-bytes")
	expectCmd := &exec.Cmd{}
	var expectStdout, expectStderr *bytes.Buffer // non-SUT

	mockExecutor := mock_exec.NewMockExecutor(mockCtrl)
	mockExecutor.EXPECT().
		Command(
			"kubectl", "config", "set-credentials", "some-user",
			"--kubeconfig", filename,
			"--token", "some-bytes",
		).
		Return(expectCmd)
	mockExecutor.EXPECT().Buffered(ctx, expectCmd).Return(expectStdout, expectStderr, cage_exec.PipelineResult{}, nil)
	client.Executor = mockExecutor

	require.NoError(t, client.UpsertUserToken(ctx, file, "some-user", expectToken))
}

func (s *ConfigSuite) TestClientUpsertContext() {
	t := s.T()
	ctx := context.Background()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	filename := filepath.Join(testkit_file.FixtureDataDir(), "kubeconfig-orig.yml")
	client := config.NewDefaultClient()

	file, err := client.Parse(filename)
	require.NoError(t, err)

	expectCmd := &exec.Cmd{}
	var expectStdout, expectStderr *bytes.Buffer // non-SUT

	mockExecutor := mock_exec.NewMockExecutor(mockCtrl)
	mockExecutor.EXPECT().
		Command(
			"kubectl", "config", "set-context", "some-context",
			"--kubeconfig", filename,
			"--cluster", "some-cluster",
			"--namespace", "some-namespace",
			"--user", "some-user",
		).
		Return(expectCmd)
	mockExecutor.EXPECT().Buffered(ctx, expectCmd).Return(expectStdout, expectStderr, cage_exec.PipelineResult{}, nil)
	client.Executor = mockExecutor

	require.NoError(t, client.UpsertContext(ctx, file, "some-context", "some-cluster", "some-namespace", "some-user"))
}
