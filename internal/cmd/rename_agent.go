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
	rename "github.com/datasance/potctl/internal/rename/agent"
	"github.com/datasance/potctl/pkg/util"
	"github.com/spf13/cobra"
)

func newRenameAgentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "agent NAME NEW_NAME",
		Short:   "Rename an Agent",
		Long:    `Rename an Agent`,
		Example: `potctl rename agent NAME NEW_NAME`,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// Get name and the new name of agent
			name := args[0]
			newName := args[1]
			namespace, err := cmd.Flags().GetString("namespace")
			util.Check(err)
			useDetached, err := cmd.Flags().GetBool("detached")
			util.Check(err)

			// Get an executor for the command
			err = rename.Execute(namespace, name, newName, useDetached)
			util.Check(err)

			util.PrintSuccess(getRenameSuccessMessage("Agent", name, newName))
		},
	}

	cmd.Flags().Bool("detached", false, pkg.flagDescDetached)

	return cmd
}
