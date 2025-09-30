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
	rebuildmicroservice "github.com/datasance/potctl/internal/rebuild/microservice"
	"github.com/datasance/potctl/pkg/util"
	"github.com/spf13/cobra"
)

func newRebuildMicroserviceCommand() *cobra.Command {
	opt := rebuildmicroservice.Options{}
	cmd := &cobra.Command{
		Use:     "microservice AppNAME/MsvcNAME",
		Short:   "Rebuilds a microservice",
		Long:    "Rebuilds a microservice",
		Example: `potctl rebuild microservice AppNAME/MsvcNAME`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) == 0 {
				util.Check(util.NewInputError("Must specify an microservice to rebuild"))
			}
			opt.Name = args[0]
			opt.Namespace, err = cmd.Flags().GetString("namespace")
			util.Check(err)

			exe := rebuildmicroservice.NewExecutor(opt)

			// Execute the command
			err = exe.Execute()
			util.Check(err)

			util.PrintSuccess("Successfully rebuild Microservice " + opt.Name)
		},
	}
	return cmd
}
