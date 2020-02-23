package testkit

import (
	"context"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	cage_k8s_config "github.com/codeactual/kubeauth/internal/cage/kubernetes/v1/kubectl/config"
)

const (
	// Prefix is prepended to most reusable data values to clarify their origin
	// when viewing them in error messages.
	Prefix = "kubeauth-testkit"

	// For explicit CLI flag selections

	ClusterName         = Prefix + "-explicit-cluster"
	ClusterRoleBindName = Prefix + "cluster-role-bind"
	ClusterRoleName     = Prefix + "cluster-role"
	ConfigFilename      = Prefix + "-kubeconfig"
	ContextName         = Prefix + "-explicit-context"
	GroupName           = Prefix + "-group"
	Namespace           = Prefix + "-explicit-namespace"
	ServiceAccountName  = Prefix + "-test-sa"

	// For missing CLI flag selections which will default to the "current" configs.

	CurrentClusterName = Prefix + "-current-cluster"
	CurrentContextName = Prefix + "-current-context"
	CurrentNamespace   = Prefix + "-current-namespace"
	Username           = Prefix + "-username"

	// Self-describing alternatives to boolean values in expected-argument lists.

	AllNamspacesEnabled  = true
	AllNamspacesDisabled = false
	Exists               = true
	NotExists            = false

	// Unsorted fixtures

	RoleBindName              = Prefix + "role-bind"
	RoleName                  = Prefix + "role"
	SecretNameSuffix          = "-token-1abcd"
	Server                    = Prefix + "-server"
	ServiceAccountSubjectName = Prefix + "-test-sa-subject"
	TokenSuffix               = "-token-1abcd"
)

// NewConfigFile returns a minimal File object which supports two modes: creation of a file
// which only has "current" values because the test case will use values of flags like --context,
// or creation of a file which contains both the "current" values and an additional context
// matching the input context/namespace value.
func NewConfigFile(filename, context, cluster, namespace string) *cage_k8s_config.File {
	f := cage_k8s_config.File{
		Name: filename,
		ClientCmdConfig: clientcmdapi.Config{
			CurrentContext: CurrentContextName,
			Clusters: map[string]*clientcmdapi.Cluster{
				CurrentClusterName: {
					Server: Server,
				},
			},
			Contexts: map[string]*clientcmdapi.Context{
				CurrentContextName: {
					AuthInfo:  Username,
					Cluster:   CurrentClusterName,
					Namespace: CurrentNamespace,
				},
			},
		},
	}

	if context != CurrentContextName {
		f.ClientCmdConfig.Contexts[context] = &clientcmdapi.Context{
			AuthInfo:  Username,
			Cluster:   cluster,
			Namespace: namespace,
		}
	}

	if cluster != CurrentClusterName {
		f.ClientCmdConfig.Clusters[cluster] = &clientcmdapi.Cluster{
			Server: Server,
		}
	}

	return &f
}

// NewNamespace returns a namespace object with common-case fields initialized.
func NewNamespace(name string) *core.Namespace {
	return &core.Namespace{ObjectMeta: meta.ObjectMeta{Name: name}}
}

// Ctx returns a common-case value for expected call parameters.
func Ctx() context.Context {
	return context.Background()
}
