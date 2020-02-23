// Copyright (C) 2020 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

//go:generate mockgen -copyright_file=$LICENSE_HEADER -package=mock -destination=$GODIR/mock/config.go -source=$GODIR/$GOFILE
package config

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	cage_exec "github.com/codeactual/kubeauth/internal/cage/os/exec"
)

// File represents a kubectl config file parsed by a Client implementation.
type File struct {
	// Name is the path to the kubectl config file.
	Name string

	// ClientCmdConfig represents the parsed config file.
	//
	// It is derived from the kubectl config file by client-go.
	ClientCmdConfig clientcmdapi.Config

	// RestConfig is the REST client config derived from clientCmdConfig.
	RestConfig *rest.Config
}

// GetCurrentContext returns the context selected in the config.
func (f *File) GetCurrentContext() (name string, _ *clientcmdapi.Context, _ error) {
	if len(f.ClientCmdConfig.Contexts) == 0 {
		return "", nil, errors.New("no context configs found")
	}

	curContext, ok := f.ClientCmdConfig.Contexts[f.ClientCmdConfig.CurrentContext]
	if curContext == nil || !ok {
		return "", nil, errors.Errorf("details of current-context [%s] not found", f.ClientCmdConfig.CurrentContext)
	}

	return f.ClientCmdConfig.CurrentContext, curContext, nil
}

// GetCurrentCluster returns the cluster selected by the current context.
func (f *File) GetCurrentCluster() (name string, _ *clientcmdapi.Cluster, _ error) {
	if len(f.ClientCmdConfig.Clusters) == 0 {
		return "", nil, errors.New("no cluster configs found")
	}

	_, curContext, err := f.GetCurrentContext()
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	cluster, ok := f.ClientCmdConfig.Clusters[curContext.Cluster]
	if cluster == nil || !ok {
		return "", nil, errors.Errorf(
			"details of current-context [%s] cluster [%s] not found",
			f.ClientCmdConfig.CurrentContext,
			curContext.Cluster,
		)
	}

	return curContext.Cluster, cluster, nil
}

// Client provides an interface to kubectl config files.
type Client interface {
	// Parse returns a Config based on the contents of the namedfile.
	Parse(filename string) (*File, error)

	// UpsertUserToken adds/updates a user's bearer token.
	UpsertUserToken(ctx context.Context, parsed *File, user string, token []byte) error

	// UpsertContext adds or updates a context.
	UpsertContext(ctx context.Context, parsed *File, name, cluster, ns, user string) error
}

// DefaultClient implementation of Client operates on real config files.
type DefaultClient struct {
	// Executor provides an os/exec.Command API for running the kubectl CLI.
	Executor cage_exec.Executor
}

func NewDefaultClient() *DefaultClient {
	return &DefaultClient{Executor: cage_exec.CommonExecutor{}}
}

// NewDefaultClient returns a DefaultClient initialized by the config file named by the
// input or, if the latter is empty, by the default from k8s.io/client-go.
func (c *DefaultClient) Parse(filename string) (*File, error) {
	file := File{}

	var loadRules *clientcmd.ClientConfigLoadingRules
	if filename == "" {
		loadRules = clientcmd.NewDefaultClientConfigLoadingRules()
		file.Name = loadRules.GetDefaultFilename()
	} else {
		loadRules = &clientcmd.ClientConfigLoadingRules{ExplicitPath: filename}
		file.Name = filename
	}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadRules,
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{}},
	)

	var err error

	file.ClientCmdConfig, err = clientConfig.RawConfig()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load config file [%s]", file.Name)
	}

	if file.ClientCmdConfig.CurrentContext == "" {
		return nil, errors.Errorf("current-context not found in file [%s]", file.Name)
	}

	file.RestConfig, err = clientConfig.ClientConfig()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create REST config from file [%s]", file.Name)
	}

	return &file, nil
}

// UpsertUserToken adds/updates a user's bearer token.
//
// It implements Client.
func (c *DefaultClient) UpsertUserToken(ctx context.Context, file *File, user string, token []byte) error {
	_, stderrBuf, _, err := c.Executor.Buffered(ctx, c.Executor.Command(
		"kubectl", "config", "set-credentials", user,
		"--kubeconfig", file.Name,
		"--token", string(token),
	))

	if err != nil {
		return errors.Wrap(err, strings.TrimSpace(stderrBuf.String()))
	}

	ctxErr := ctx.Err()
	if ctxErr != nil {
		return errors.WithStack(ctxErr)
	}

	return nil
}

// UpsertContext adds or updates a context.
//
// It implements Client.
func (c *DefaultClient) UpsertContext(ctx context.Context, file *File, name, cluster, ns, user string) error {
	_, stderrBuf, _, err := c.Executor.Buffered(ctx, c.Executor.Command(
		"kubectl", "config", "set-context", name,
		"--kubeconfig", file.Name,
		"--cluster", cluster,
		"--namespace", ns,
		"--user", user,
	))

	if err != nil {
		return errors.Wrap(err, strings.TrimSpace(stderrBuf.String()))
	}

	ctxErr := ctx.Err()
	if ctxErr != nil {
		return errors.WithStack(ctxErr)
	}

	return nil
}

var _ Client = (*DefaultClient)(nil)
