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
	delete "github.com/datasance/potctl/internal/delete/volume"
	"github.com/datasance/potctl/pkg/util"
	"github.com/spf13/cobra"
)

func newDeleteVolumeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume NAME",
		Short: "Delete an Volume",
		Long: `Delete an Volume.

The Volume will be deleted from the Agents that it is stored on.`,
		Example: `potctl delete volume NAME`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Get name and namespace of volume
			name := args[0]
			namespace, err := cmd.Flags().GetString("namespace")
			util.Check(err)

			// Run the command
			exe, err := delete.NewExecutor(namespace, name)
			util.Check(err)
			err = exe.Execute()
			util.Check(err)

			util.PrintSuccess("Successfully deleted " + namespace + "/" + name)
		},
	}

	return cmd
}
