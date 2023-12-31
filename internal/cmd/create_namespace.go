/*
 *  *******************************************************************************
 *  * Copyright (c) 2023 Datasance Teknoloji A.S.
 *  *
 *  * This program and the accompanying materials are made available under the
 *  * terms of the Eclipse Public License v. 2.0 which is available at
 *  * http://www.eclipse.org/legal/epl-2.0
 *  *
 *  * SPDX-License-Identifier: EPL-2.0
 *  *******************************************************************************
 *
 */

package cmd

import (
	create "github.com/datasance/potctl/internal/create/namespace"
	"github.com/datasance/potctl/pkg/util"
	"github.com/spf13/cobra"
)

func newCreateNamespaceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "namespace NAME",
		Short: "Create a Namespace",
		Long: `Create a Namespace.

A Namespace contains all components of an Edge Compute Network.

A single instance of potctl can be used to manage any number of Edge Compute Networks.`,
		Example: `potctl create namespace NAME`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Get name and namespace of agent
			name := args[0]

			// Run the command
			err := create.Execute(name)
			util.Check(err)

			util.PrintSuccess("Successfully created namespace " + name)
		},
	}

	return cmd
}
