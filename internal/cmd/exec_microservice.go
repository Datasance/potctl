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
	"github.com/datasance/potctl/internal/exec"
	"github.com/datasance/potctl/pkg/util"
	"github.com/spf13/cobra"
)

func newExecMicroserviceCommand() *cobra.Command {
	opt := exec.Options{
		Resource: "microservice",
	}

	cmd := &cobra.Command{
		Use:     "microservice AppName/MsvcName",
		Short:   "Connect to an Exec Session of a Microservice",
		Long:    `Connect to an Exec Session of a Microservice to interact with its container.`,
		Example: `potctl exec microservice AppName/MicroserviceName`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Get resource type and name
			var err error
			opt.Name = args[0]
			opt.Namespace, err = cmd.Flags().GetString("namespace")
			util.Check(err)

			// Get executor for exec command
			exe, err := exec.NewExecutor(&opt)
			util.Check(err)

			// Execute the command
			err = exe.Execute()
			util.Check(err)
		},
	}

	return cmd
}
