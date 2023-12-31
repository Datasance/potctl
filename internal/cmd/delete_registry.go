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
	delete "github.com/datasance/potctl/internal/delete/registry"
	"github.com/datasance/potctl/pkg/util"
	"github.com/spf13/cobra"
)

func newDeleteRegistryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "registry ID",
		Short:   "Delete a Registry",
		Long:    `Delete a Registry from the Controller.`,
		Example: `potctl delete registry ID`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Get name and namespace
			id := args[0]
			namespace, err := cmd.Flags().GetString("namespace")
			util.Check(err)

			// Get an executor for the command
			exe, err := delete.NewExecutor(namespace, id)
			util.Check(err)
			err = exe.Execute()
			util.Check(err)

			util.PrintSuccess("Successfully deleted registry " + id)
		},
	}

	return cmd
}
