// Copyright (C) 2020 The kubeauth Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package add_user

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"

	"github.com/codeactual/kubeauth/internal/cage/cli/handler"
	handler_cobra "github.com/codeactual/kubeauth/internal/cage/cli/handler/cobra"
	cage_k8s "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1"
	cage_k8s_core "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/core"
	cage_k8s_config "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/kubectl/config"
	cage_k8s_rbac "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac"
	cage_file "github.com/codeactual/kubeauth/internal/cage/os/file"
	cage_reflect "github.com/codeactual/kubeauth/internal/cage/reflect"
)

// Handler defines the sub-command flags and logic.
type Handler struct {
	handler.Session

	KubeApiClientset    *cage_k8s_core.Clientset
	KubectlConfigClient cage_k8s_config.Client

	Cluster            string   `usage:"cluster of the new context to create (default from current-context)"`
	ClusterRoles       []string `usage:"cluster role binding to create (<role name>:<binding name>)"`
	ConfigFile         string   `usage:"kubectl config file to modify"`
	Namespace          string   `usage:"namespace to receive service account (default from current-context)"`
	Roles              []string `usage:"role binding to create (<role name>:<binding name>)"`
	ServiceAccountName string   `usage:"name of service account to create"`
	Username           string   `usage:"username/context to receive the service account's bearer token"`

	// Verbosity levels greater than 0 will enable status messages and error stack traces.
	//
	// It is an int for consistency with other commands, even though levels beyond 1 are not used.
	Verbosity int `usage:"kubectl verbosity level"`
}

// Init defines the command, its environment variable prefix, etc.
//
// It implements cli/handler/cobra.Handler.
func (h *Handler) Init() handler_cobra.Init {
	return handler_cobra.Init{
		Cmd: &cobra.Command{
			Use:   "add-user",
			Short: "Add a service account and user/context to use its credentials",
		},
		EnvPrefix: "KUBEAUTH",
	}
}

// BindFlags binds the flags to Handler fields.
//
// It implements cli/handler/cobra.Handler.
func (h *Handler) BindFlags(cmd *cobra.Command) []string {
	cmd.Flags().StringVarP(&h.Cluster, "cluster", "", "", cage_reflect.GetFieldTag(*h, "Cluster", "usage"))
	cmd.Flags().StringSliceVarP(&h.ClusterRoles, "cluster-role", "", []string{}, cage_reflect.GetFieldTag(*h, "ClusterRoles", "usage"))
	cmd.Flags().StringVarP(&h.ConfigFile, "kubeconfig", "", "", cage_reflect.GetFieldTag(*h, "ConfigFile", "usage"))
	cmd.Flags().StringVarP(&h.Namespace, "namespace", "n", "", cage_reflect.GetFieldTag(*h, "Namespace", "usage"))
	cmd.Flags().StringVarP(&h.ServiceAccountName, "account", "", "", cage_reflect.GetFieldTag(*h, "ServiceAccountName", "usage"))
	cmd.Flags().StringSliceVarP(&h.Roles, "role", "", []string{}, cage_reflect.GetFieldTag(*h, "Roles", "usage"))
	cmd.Flags().StringVarP(&h.Username, "user", "", "", cage_reflect.GetFieldTag(*h, "Username", "usage"))
	cmd.Flags().IntVarP(&h.Verbosity, "v", "v", 0, cage_reflect.GetFieldTag(*h, "Verbosity", "usage"))
	return []string{"account", "user"}
}

// Run performs the sub-command logic.
//
// It implements cli/handler/cobra.Handler.
//
// Based on:
//   https://stackoverflow.com/questions/42170380/how-to-add-users-to-kubernetes-kubectl/42186135#42186135
//   https://gist.github.com/innovia/fbba8259042f71db98ea8d4ad19bd708
func (h *Handler) Run(ctx context.Context, input handler.Input) {
	if err := h.run(ctx, input); err != nil {
		if h.Verbosity > 0 {
			h.ExitOnErr(err, "", 1)
		} else {
			h.ExitOnErrShort(err, "", 1)
		}
	}
}

func (h *Handler) run(ctx context.Context, _ handler.Input) error {
	stderr := h.Err()
	verbose := func(format string, vArgs ...interface{}) {
		if h.Verbosity > 0 {
			fmt.Fprintln(stderr, "kubeauth: "+fmt.Sprintf(format, vArgs...))
		}
	}

	var roleBindings, clusterRoleBindings []*cage_k8s_rbac.BindingSelector

	// Create clients.

	configClient := h.KubectlConfigClient
	if configClient == nil {
		configClient = cage_k8s_config.NewDefaultClient()
	}

	configFile, err := configClient.Parse(h.ConfigFile)
	if err != nil {
		return errors.Wrap(err, "kubeauth") // WithStack alternative to disambiguate from kubectl output
	}

	apiClientset := h.KubeApiClientset
	if apiClientset == nil {
		rawApiClientset, err := kubernetes.NewForConfig(configFile.RestConfig)
		if err != nil {
			return errors.Wrap(err, "kubeauth: failed to create API client")
		}

		apiClientset = cage_k8s_core.NewClientset(rawApiClientset)
	}

	roleClient := apiClientset.Roles
	clusterRoleClient := apiClientset.ClusterRoles
	roleBindingClient := apiClientset.RoleBindings
	clusterRoleBindingClient := apiClientset.ClusterRoleBindings
	secretClient := apiClientset.Secrets
	saClient := apiClientset.ServiceAccounts

	// Validate inputs.

	if h.Cluster == "" {
		h.Cluster, _, err = configFile.GetCurrentCluster()
		if err != nil {
			return errors.Wrap(err, "kubeauth") // WithStack alternative to disambiguate from kubectl output
		}
	}

	// - Mirror the behavior of kubectl regarding --namespace and the current context.
	if h.Namespace == "" {
		_, curContext, err := configFile.GetCurrentContext()
		if err != nil {
			return errors.Wrap(err, "kubeauth") // WithStack alternative to disambiguate from kubectl output
		}

		h.Namespace = curContext.Namespace
	}

	if len(h.Roles) > 0 {
		var invalid []string
		for _, r := range h.Roles {
			binding, err := cage_k8s_rbac.NewBindingSelector(r)
			if err != nil {
				invalid = append(invalid, err.Error())
				continue
			}
			roleBindings = append(roleBindings, binding)
		}
		if len(invalid) > 0 {
			return errors.Errorf("kubeauth: invalid --role selectors: %q", invalid)
		}

		invalid = []string{}
		for _, b := range roleBindings {
			_, exists, err := roleClient.Get(h.Namespace, b.RoleName)
			if err != nil {
				return errors.Wrap(err, "kubeauth") // WithStack alternative to disambiguate from kubectl output
			}

			if !exists {
				invalid = append(invalid, b.RoleName)
			}
		}
		if len(invalid) > 0 {
			return errors.Errorf("kubeauth: role(s) not found: %q", invalid)
		}
	}

	if len(h.ClusterRoles) > 0 {
		var invalid []string
		for _, r := range h.ClusterRoles {
			binding, err := cage_k8s_rbac.NewBindingSelector(r)
			if err != nil {
				invalid = append(invalid, err.Error())
				continue
			}
			clusterRoleBindings = append(clusterRoleBindings, binding)
		}
		if len(invalid) > 0 {
			return errors.Errorf("kubeauth: invalid --cluster-role selectors: %q", invalid)
		}

		invalid = []string{}
		for _, b := range clusterRoleBindings {
			_, exists, err := clusterRoleClient.Get(b.RoleName)
			if err != nil {
				return errors.Wrap(err, "kubeauth") // WithStack alternative to disambiguate from kubectl output
			}

			if !exists {
				invalid = append(invalid, b.RoleName)
			}
		}
		if len(invalid) > 0 {
			return errors.Errorf("kubeauth: cluster role(s) not found: %q", invalid)
		}
	}

	// Create the service account if needed.

	var saObj *core.ServiceAccount
	var exists bool

	saObj, exists, err = saClient.Get(h.Namespace, h.ServiceAccountName)
	if err != nil {
		return errors.Wrap(err, "kubeauth") // WithStack alternative to disambiguate from kubectl output
	}

	if exists {
		verbose("service account already exists")
	} else {
		saObj, err = saClient.CreateBasic(h.Namespace, h.ServiceAccountName)
		if err != nil {
			return errors.Wrap(err, "kubeauth") // WithStack alternative to disambiguate from kubectl output
		}
	}

	// Get the name of the secret which holds the service account's token.
	// Use polling to work around the fact that the API call to create a service
	// account may return successfully without having created the token. It is due to
	// how the SA token controller is responsible for creating the token and completes
	// its work independently.
	//
	// https://github.com/kubernetes/kubernetes/blob/v1.17.0/pkg/controller/serviceaccount/tokens_controller.go#L358
	//
	// Kubernetes uses polling in its own coverage of token creation:
	// https://github.com/kubernetes/kubernetes/blob/v1.17.0/test/integration/serviceaccount/service_account_test.go#L124

	// Assumes that (at least in the v1 API) the auto-generated token is stored as a
	// secret with this naming convention: "<account name>-token-<random>".
	saTokenAvail := func() bool {
		for _, s := range saObj.Secrets {
			if strings.HasPrefix(s.Name, h.ServiceAccountName+"-token-") {
				return true
			}
		}
		return false
	}

	if !saTokenAvail() {
		backoffOpts := wait.Backoff{
			Duration: 100 * time.Millisecond,
			Factor:   2,
			// The sleep at each iteration is the duration plus an additional
			// amount chosen uniformly at random from the interval between
			// zero and `jitter*duration`.
			Jitter: 1,
			Steps:  5,
		}
		backoffCond := func() (done bool, err error) {
			saObj, _, err = saClient.Get(h.Namespace, h.ServiceAccountName)
			if err != nil {
				return false, errors.WithStack(err)
			}
			return saTokenAvail(), nil
		}

		err = wait.ExponentialBackoff(backoffOpts, backoffCond)
		if err != nil {
			return errors.Wrap(err, "kubeauth: secret with service account's token not found")
		}
	}

	// Bind the service account to selected roles, if any.

	for _, b := range roleBindings {
		_, err = roleBindingClient.Create(
			h.Namespace, b.BindingName, b.RoleName,
			rbac.Subject{Kind: cage_k8s.KindServiceAccount, Name: h.ServiceAccountName},
		)
		if err != nil {
			if k8s_errors.IsAlreadyExists(err) {
				verbose("role binding(s) already exist")
				continue
			}
			return errors.Wrap(err, "kubeauth") // WithStack alternative to disambiguate from kubectl output
		}
	}

	for _, b := range clusterRoleBindings {
		_, err = clusterRoleBindingClient.Create(
			b.BindingName, b.RoleName,
			rbac.Subject{Namespace: h.Namespace, Kind: cage_k8s.KindServiceAccount, Name: h.ServiceAccountName},
		)
		if err != nil {
			if k8s_errors.IsAlreadyExists(err) {
				verbose("cluster role binding(s) already exist")
				continue
			}
			return errors.Wrap(err, "kubeauth") // WithStack alternative to disambiguate from kubectl output
		}
	}

	// Retrieve and write the service account's token to a temporary file.

	secretObj, exists, err := secretClient.Get(h.Namespace, saObj.Secrets[0].Name)
	if err != nil {
		return errors.Wrapf(err, "kubeauth: failed to query for service account's secret [%s]", saObj.Secrets[0].Name)
	}

	if !exists {
		return errors.Errorf("kubeauth: service account's secret [%s] not found", saObj.Secrets[0].Name)
	}

	caCrt, ok := secretObj.Data["ca.crt"]
	if !ok {
		return errors.Errorf("kubeauth: service account's secret [%s] does not contain 'ca.crt' data", saObj.Secrets[0].Name)
	}

	token, ok := secretObj.Data["token"]
	if !ok {
		return errors.Errorf("kubeauth: service account's secret [%s] does not contain 'token' data", saObj.Secrets[0].Name)
	}

	caCrtFile, err := ioutil.TempFile("", "kubeauth.*.ca.crt")
	if err != nil {
		return errors.Wrap(err, "kubeuath: failed to create temporary ca.crt file")
	}

	caCrtFilename := caCrtFile.Name()
	defer cage_file.RemoveSafer(caCrtFilename)

	if _, err = caCrtFile.Write(caCrt); err != nil {
		return errors.Wrap(err, "kubeauth: failed to write temporary ca.crt file")
	}
	if err = caCrtFile.Sync(); err != nil {
		return errors.Wrap(err, "kubeauth: failed to sync temporary ca.crt file")
	}
	if err = caCrtFile.Close(); err != nil {
		return errors.Wrap(err, "kubeauth: failed to close temporary ca.crt file")
	}

	// Add/update a user in the config file which authenticates using the service account's token.

	if err = configClient.UpsertUserToken(ctx, configFile, h.Username, token); err != nil {
		return errors.Wrap(err, "kubeauth: failed to set user token")
	}

	// Name the context after the username.
	if err = configClient.UpsertContext(ctx, configFile, h.Username, h.Cluster, h.Namespace, h.Username); err != nil {
		errors.Wrap(err, "kubeauth: failed to set user token")
	}

	return nil
}

// New returns a cobra command instance based on Handler.
func NewCommand() *cobra.Command {
	return handler_cobra.NewHandler(&Handler{
		Session: &handler.DefaultSession{},
	})
}

var _ handler_cobra.Handler = (*Handler)(nil)
