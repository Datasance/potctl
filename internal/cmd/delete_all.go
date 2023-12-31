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
	delete "github.com/datasance/potctl/internal/delete/all"
	"github.com/datasance/potctl/pkg/util"
	"github.com/spf13/cobra"
)

func newDeleteAllCommand() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "all",
		Short: "Delete all resources within a namespace",
		Long: `Delete all resources within a namespace.

Tears down all components of an Edge Compute Network.

If you don't want to tear down the deployments but would like to free up the Namespace, use the disconnect command instead.`,
		Example: `potctl delete all -n NAMESPACE`,
		Run: func(cmd *cobra.Command, args []string) {
			// Execute command
			namespace, err := cmd.Flags().GetString("namespace")
			util.Check(err)
			useDetached, err := cmd.Flags().GetBool("detached")
			util.Check(err)
			err = delete.Execute(namespace, useDetached, force)
			util.Check(err)

			util.PrintSuccess("Successfully deleted all resources in namespace " + namespace)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Force deletion of Agents")
	cmd.Flags().Bool("detached", false, pkg.flagDescDetached)

	return cmd
}
