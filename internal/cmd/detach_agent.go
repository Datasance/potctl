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
	detach "github.com/datasance/potctl/internal/detach/agent"
	"github.com/datasance/potctl/pkg/util"
	"github.com/spf13/cobra"
)

func newDetachAgentCommand() *cobra.Command {
	force := false
	cmd := &cobra.Command{
		Use:   "agent NAME",
		Short: "Detaches an Agent",
		Long: `Detaches an Agent.

The Agent will be deprovisioned from the Controller within the namespace.
The Agent will be removed from Controller.

You cannot detach unprovisioned Agents.

The Agent stack will not be uninstalled from the host.`,
		Example: `potctl detach agent NAME`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Get name and namespace of agent
			name := args[0]
			namespace, err := cmd.Flags().GetString("namespace")
			util.Check(err)

			// Run the command
			exe := detach.NewExecutor(namespace, name, force)
			err = exe.Execute()
			util.Check(err)

			util.PrintSuccess("Successfully detached " + name)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Detach Agent even if it is running Microservices")

	return cmd
}
