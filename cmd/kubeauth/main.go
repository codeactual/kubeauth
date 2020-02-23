// Copyright (C) 2020 The kubeauth Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Command kubeauth assists with authentication-related maintenance tasks.
//
// Usage:
//
//   kubeauth --help
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/codeactual/kubeauth/cmd/kubeauth/add_user"
	"github.com/codeactual/kubeauth/cmd/kubeauth/ctl"
	"github.com/codeactual/kubeauth/internal/cage/cli/handler"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "kubeauth",
	}

	rootCmd.Version = handler.Version()
	rootCmd.AddCommand(add_user.NewCommand())
	rootCmd.AddCommand(ctl.NewCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n%+v\n", rootCmd.UsageString(), err)
	}
}
