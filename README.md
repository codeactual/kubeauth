# kubeauth [![GoDoc](https://godoc.org/github.com/codeactual/kubeauth?status.svg)](https://pkg.go.dev/mod/github.com/codeactual/kubeauth) [![Go Report Card](https://goreportcard.com/badge/github.com/codeactual/kubeauth)](https://goreportcard.com/report/github.com/codeactual/kubeauth) [![Build Status](https://travis-ci.org/codeactual/kubeauth.png)](https://travis-ci.org/codeactual/kubeauth)

kubeauth is a program to assist usage of `kubectl` for user/group related operations. It currently provides two commands:

1. `add-user` creates a service account based user, adds the credentials to the selected kubeconfig, and optionally creates bindings to existing roles or cluster roles.
1. `ctl` wraps `kubectl` invocation and validates flags such as `--as` and `--as-group`.

## `add-user`

### Examples

> Create the kubeconfig user "tester" based on service account "default" in the "dev" namespace. Also bind it to a role and cluster role. The --role and --cluster-role flags may be supplied multiple times.

```bash
kubeauth add-user -v=1 \
  --user tester \
  --account default \
  --namespace dev \
  --role role_name_0:binding_name_0 \
  --cluster-role role_name_1:binding_name_1
```

### Validation checks

- `--role`: role exists in effective namespace
- `--cluster-role`: cluster role exists

## `ctl`

- Invocation format: `ctl [kubectl sub-command] [kubeauth flags] -- [kubectl sub-command flags]`
- `ctl` flags which are also accepted by `kubectl` will be passed to the latter.

### Examples

> Verify that "tester" exists and run "kubectl auth can-i -v=1 --as tester --list".

```bash
kubeauth ctl auth can-i -v=1 \
  --as tester \
  -- --list
```

> Verify that "system:serviceaccount:dev:default" exists and run "kubectl auth can-i -v=1 --as system:serviceaccount:dev:default --list".

```bash
kubeauth ctl auth can-i -v=1 \
  --as system:serviceaccount:dev:default \
  -- --list
```

### Validation checks

- effective context exists
- effective namespace exists
- `--as` selection exists
- `--as-group` selection exists
- agreement between `--cluster` and effective context's cluster

# Development

## License

[Mozilla Public License Version 2.0](https://www.mozilla.org/en-US/MPL/2.0/) ([About](https://www.mozilla.org/en-US/MPL/), [FAQ](https://www.mozilla.org/en-US/MPL/2.0/FAQ/))

- `add_user` was based on [this bash script gist](https://gist.github.com/innovia/fbba8259042f71db98ea8d4ad19bd708).

## Contributing

- Please feel free to submit issues, PRs, questions, and feedback.
- Although this repository consists of snapshots extracted from a private monorepo using [transplant](https://github.com/codeactual/transplant), PRs are welcome. Standard GitHub workflows are still used.

## Testing

### `ctl`

- Reminders
  - "you typically need to include `--as-group=system:authenticated` in order to have permission to run a `selfsubjectaccessreview` check." (https://github.com/kubernetes/kubernetes/issues/73123#issuecomment-456185028)

# FAQ

- `ctl`
  - Q: When verbose output is enabled with `-v=1` and I use `--as`/`--as-group` flags, why do I not always see `in namespace X` in the messages describing where the user/group was found?
    - A: It may be that the `--as/--as-group` identity was found in a role or cluster-role binding where the `Subject` object contained an empty `Namespace` field. At the time this was written, the empty value is expected for `User` and `Group` subjects because those object kinds are considered ["non-namespace"](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#subject-v1-rbac-authorization-k8s-io).
