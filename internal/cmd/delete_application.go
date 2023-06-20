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
	delete "github.com/datasance/potctl/v1/internal/delete/application"
	"github.com/datasance/potctl/v1/pkg/util"
	"github.com/spf13/cobra"
)

func newDeleteApplicationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "application NAME",
		Short:   "Delete an application",
		Long:    `Delete an application and all its components`,
		Example: `potctl delete application NAME`,
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Get microservice name
			name := args[0]
			namespace, err := cmd.Flags().GetString("namespace")
			util.Check(err)

			// Execute command
			err = delete.Execute(namespace, name)
			util.Check(err)

			util.PrintSuccess("Successfully deleted " + namespace + "/" + name)
		},
	}

	return cmd
}
