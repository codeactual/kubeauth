// Copyright (C) 2020 The kubeauth Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctl

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/codeactual/kubeauth/internal/cage/cli/handler"
	handler_cobra "github.com/codeactual/kubeauth/internal/cage/cli/handler/cobra"
	cage_k8s "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1"
	cage_k8s_core "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/core"
	cage_k8s_config "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/kubectl/config"
	cage_k8s_identity "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/rbac/identity"
	cage_exec "github.com/codeactual/kubeauth/internal/cage/os/exec"
	cage_reflect "github.com/codeactual/kubeauth/internal/cage/reflect"
)

// Handler defines the sub-command flags and logic.
type Handler struct {
	handler.Session

	KubectlConfigClient cage_k8s_config.Client
	KubeApiClientset    *cage_k8s_core.Clientset
	IdentityRegistry    *cage_k8s_identity.Registry

	// Executor provides an os/exec.Command API for running the kubectl CLI.
	Executor cage_exec.Executor

	AllNamespaces bool     `usage:"include identities from any/no namespace"`
	As            string   `usage:"User/ServiceAccount/Role/ClusterRole to impersonate"`
	AsGroup       []string `usage:"Group(s) to impersonate"`
	Cluster       string   `usage:"pass to kubctl if effective context's cluster matches, else error (default from current-context)"`
	ConfigFile    string   `usage:"kubectl config file to modify"`
	Context       string   `usage:"consider users in this --kubeconfig context (defaults to current-context)"`
	Namespace     string   `usage:"include identities from only one namespace (default from --context)"`

	Verbosity int `usage:"kubectl verbosity level (and verbose kubeauth output for any level > 0)"`

	// usage is the auto-generated flag-usage content.
	usage string
}

// Init defines the command, its environment variable prefix, etc.
//
// It implements cli/handler/cobra.Handler.
func (h *Handler) Init() handler_cobra.Init {
	return handler_cobra.Init{
		Cmd: &cobra.Command{
			Use:   "ctl",
			Short: "Run 'kubectl' with additional input validation",
			Long:  "kubeauth ctl [kubectl sub-command] [kubeauth flags] -- [kubectl sub-command flags]",
		},
		EnvPrefix: "KUBEAUTH",
	}
}

// BindFlags binds the flags to Handler fields.
//
// It implements cli/handler/cobra.Handler.
func (h *Handler) BindFlags(cmd *cobra.Command) []string {
	// Reminder: Use a unique flag name unless kubeauth's behavior is consistent with kubectl's
	// when the flag is passed through to the latter.
	//
	// For example, "kube-user" is used here instead of "user" because kubectl interprets it as
	// "The name of the kubeconfig user to use".
	//
	// If the flag is shared between the tools, mirror both the long and short forms.
	cmd.Flags().StringVarP(&h.ConfigFile, "kubeconfig", "", "", cage_reflect.GetFieldTag(*h, "ConfigFile", "usage"))
	cmd.Flags().StringVarP(&h.Context, "context", "", "", cage_reflect.GetFieldTag(*h, "Context", "usage"))
	cmd.Flags().StringVarP(&h.Cluster, "cluster", "", "", cage_reflect.GetFieldTag(*h, "Cluster", "usage"))
	cmd.Flags().StringVarP(&h.As, "as", "", "", cage_reflect.GetFieldTag(*h, "As", "usage"))
	cmd.Flags().StringSliceVarP(&h.AsGroup, "as-group", "", []string{}, cage_reflect.GetFieldTag(*h, "AsGroup", "usage"))
	cmd.Flags().StringVarP(&h.Namespace, "namespace", "n", "", cage_reflect.GetFieldTag(*h, "Namespace", "usage"))
	cmd.Flags().BoolVarP(&h.AllNamespaces, "all-namespaces", "", false, cage_reflect.GetFieldTag(*h, "AllNamespaces", "usage"))
	cmd.Flags().IntVarP(&h.Verbosity, "v", "v", 0, cage_reflect.GetFieldTag(*h, "Verbosity", "usage"))

	h.usage = cmd.UsageString()

	return []string{}
}

// Run performs the sub-command logic.
//
// It implements cli/handler/cobra.Handler.
func (h *Handler) Run(ctx context.Context, input handler.Input) {
	if err := h.run(ctx, input); err != nil {
		if h.Verbosity > 0 {
			h.ExitOnErr(err, "", 1)
		} else {
			h.ExitOnErrShort(err, "", 1)
		}
	}
}

func (h *Handler) run(ctx context.Context, input handler.Input) error {
	stderr := h.Err()
	verbose := func(format string, vArgs ...interface{}) {
		if h.Verbosity > 0 {
			fmt.Fprintln(stderr, "kubeauth: "+fmt.Sprintf(format, vArgs...))
		}
	}

	if h.Executor == nil {
		h.Executor = cage_exec.CommonExecutor{}
	}

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

	regClient := h.IdentityRegistry
	if regClient == nil {
		regClient = cage_k8s_identity.NewRegistry(apiClientset)
	}

	nsClient := apiClientset.Namespaces

	// Validate inputs.

	if h.As == "" && len(h.AsGroup) == 0 {
		return errors.Errorf("kubeauth: %s\nmissing --as or --as-group selection", h.usage)
	}
	if h.Namespace != "" && h.AllNamespaces {
		return errors.Errorf("kubeauth: %s\n--namespace and --all-namespaces cannot be combined", h.usage)
	}

	currentContextName, currentContext, err := configFile.GetCurrentContext()
	if err != nil {
		return errors.Wrap(err, "kubeauth") // WithStack alternative to disambiguate from kubectl output
	}

	// These describe the effective context, after considering --context and current-context values.
	var effectiveContext *clientcmdapi.Context
	var effectiveContextName string

	if h.Context == "" {
		effectiveContextName = currentContextName

		verbose(
			"defaulting to current-context [%s] from file [%s]",
			effectiveContextName, configFile.Name,
		)
	} else {
		effectiveContextName = h.Context

		verbose(
			"using --context [%s] from file [%s]",
			effectiveContextName, configFile.Name,
		)
	}

	if effectiveContext = configFile.ClientCmdConfig.Contexts[effectiveContextName]; effectiveContext == nil {
		return errors.Errorf("kubeauth: context [%s] not found in config [%s]", effectiveContextName, configFile.Name)
	}

	// Just as impersonation targets are validated, also catch cluster mismatches before passing
	// --cluster on to kubectl.
	if h.Cluster != "" && h.Cluster != effectiveContext.Cluster {
		return errors.Errorf("kubeauth: selected --cluster [%s] differs from effective context's cluster [%s]", h.Cluster, effectiveContext.Cluster)
	}

	if !h.AllNamespaces {
		// - Mirror the behavior of kubectl regarding --namespace and the current context.
		// - Validate the effective namespace selection.
		if h.Namespace == "" {
			h.Namespace = effectiveContext.Namespace

			verbose(
				"defaulting to namespace [%s] from context [%s]",
				h.Namespace, effectiveContextName,
			)
		} else {
			verbose("using --namespace [%s]", h.Namespace)
		}

		if h.Namespace != "" {
			_, exists, err := nsClient.Get(h.Namespace)
			if err != nil {
				return errors.Wrap(err, "kubeauth: failed to validate namespace")
			}

			if !exists {
				return errors.Errorf("kubeauth: selected --namespace [%s] not found", h.Namespace)
			}
		}
	}

	// Validate the impersonation target.

	// For --as, limit search scope to User identities.
	// For --as-group, limit search scope to Group identities.
	// Otherwise include all object kinds in results.
	//
	// Each QueryOption passed to Query expands the kind-scope of queriers used and/or
	// potentially further limits the results.
	if h.As != "" {
		verbose("validating --as [%s]", h.As)

		list, err := regClient.Query(
			ctx,
			cage_k8s_identity.QueryKind(cage_k8s.KindUser),
			cage_k8s_identity.QueryNamespace(h.Namespace),
			cage_k8s_identity.QueryName(h.As),
			cage_k8s_identity.QueryClientCmdConfig(&configFile.ClientCmdConfig),
		)
		if err != nil {
			return errors.Wrap(err, "kubeauth: query did not complete")
		}

		if len(list.Items) == 0 {
			return errors.Errorf("kubeauth: --as identity [%s] not found", h.As)
		}

		for _, item := range list.Items {
			verbose("--as identity found in [%s]", item)
		}
	}

	if len(h.AsGroup) > 0 {
		verbose("validating --as-group %v", h.AsGroup)

		for _, group := range h.AsGroup {
			list, err := regClient.Query(
				ctx,
				cage_k8s_identity.QueryKind(cage_k8s.KindGroup),
				cage_k8s_identity.QueryNamespace(h.Namespace),
				cage_k8s_identity.QueryName(group),
			)
			if err != nil {
				return errors.Wrap(err, "kubeauth: query did not complete")
			}

			if len(list.Items) == 0 {
				return errors.Errorf("kubeauth: --as-group identity [%s] not found", group)
			}

			for _, item := range list.Items {
				verbose("--as-group [%s] identity found in [%s]", group, item)
			}
		}
	}

	// Passthrough validated inputs to kubectl.

	useCurrentContext := effectiveContextName == currentContextName

	kubectlArgs := input.ArgsBeforeDash
	kubectlArgs = append(kubectlArgs, "--kubeconfig", configFile.Name)
	if h.As != "" {
		kubectlArgs = append(kubectlArgs, "--as", h.As)
	}
	if len(h.AsGroup) > 0 {
		for _, group := range h.AsGroup {
			kubectlArgs = append(kubectlArgs, "--as-group", group)
		}
	}
	if !useCurrentContext {
		kubectlArgs = append(kubectlArgs, "--context", effectiveContextName)
	}
	if h.Cluster != "" && (!useCurrentContext || h.Cluster != currentContext.Cluster) {
		kubectlArgs = append(kubectlArgs, "--cluster", h.Cluster)
	}
	if h.AllNamespaces {
		kubectlArgs = append(kubectlArgs, "--all-namespaces")
	} else if h.Namespace != "" && h.Namespace != currentContext.Namespace {
		kubectlArgs = append(kubectlArgs, "--namespace", h.Namespace)
	}
	kubectlArgs = append(kubectlArgs, "--v", strconv.Itoa(h.Verbosity))

	kubectlArgs = append(kubectlArgs, input.ArgsAfterDash...)
	kubectlCmd := h.Executor.Command("kubectl", kubectlArgs...)

	verbose("running: kubectl %s", strings.Join(kubectlArgs, " "))

	res, err := h.Executor.Standard(ctx, h.Out(), h.Err(), nil, kubectlCmd)
	if err != nil {
		os.Exit(res.Cmd[kubectlCmd].Code)
	}

	ctxErr := ctx.Err()
	if ctxErr != nil {
		return errors.Wrap(err, "kubeauth: command cancelled")
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
