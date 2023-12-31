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
	delete "github.com/datasance/potctl/internal/delete/namespace"
	"github.com/datasance/potctl/pkg/util"
	"github.com/spf13/cobra"
)

func newDeleteNamespaceCommand() *cobra.Command {
	force := false
	cmd := &cobra.Command{
		Use:   "namespace NAME",
		Short: "Delete a Namespace",
		Long: `Delete a Namespace.

The Namespace must be empty.

If you would like to delete all resources in the Namespace, use the --force flag.`,
		Example: `potctl delete namespace NAME`,
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Get microservice name
			name := args[0]

			// Execute command
			err := delete.Execute(name, force)
			util.Check(err)

			util.PrintSuccess("Successfully deleted Namespace " + name)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Force deletion of all resources within the Namespace")

	return cmd
}
