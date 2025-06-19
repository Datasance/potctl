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
	"fmt"
	detach "github.com/datasance/potctl/internal/detach/volumemount"
	"github.com/datasance/potctl/pkg/util"
	"github.com/spf13/cobra"
	"strings"
)

func newDetachVolumeMountCommand() *cobra.Command {
	opt := detach.Options{}
	cmd := &cobra.Command{
		Use:     "volume-mount NAME AGENT_NAME1 AGENT_NAME2",
		Short:   "Detach a Volume Mount from existing Agents",
		Long:    `Detach a Volume Mount from existing Agents.`,
		Example: `potctl detach volume-mount NAME AGENT_NAME1 AGENT_NAME2`,
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// Get name and namespace of agent
			opt.Name = args[0]
			opt.Agents = args[1:]
			var err error
			opt.Namespace, err = cmd.Flags().GetString("namespace")
			util.Check(err)

			// Run the command
			exe := detach.NewExecutor(opt)
			err = exe.Execute()
			util.Check(err)

			msg := fmt.Sprintf("Successfully detached Volume Mount %s from Agents %s", opt.Name, strings.Join(opt.Agents, ", "))
			util.PrintSuccess(msg)
		},
	}

	return cmd
}
