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
	deleteserviceaccount "github.com/datasance/potctl/internal/delete/serviceaccount"
	"github.com/datasance/potctl/pkg/util"
	"github.com/spf13/cobra"
)

func newDeleteServiceAccountCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "serviceaccount NAME",
		Short:   "Delete a ServiceAccount",
		Long:    `Delete a ServiceAccount from the Controller.`,
		Example: `potctl delete serviceaccount NAME`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			namespace, err := cmd.Flags().GetString("namespace")
			util.Check(err)

			exe, err := deleteserviceaccount.NewExecutor(namespace, name)
			util.Check(err)
			err = exe.Execute()
			util.Check(err)

			util.PrintSuccess("Successfully deleted serviceaccount " + name)
		},
	}

	return cmd
}
