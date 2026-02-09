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
	"github.com/datasance/potctl/internal/describe"
	"github.com/datasance/potctl/pkg/util"
	"github.com/spf13/cobra"
)

func newDescribeServiceAccountCommand() *cobra.Command {
	opt := describe.Options{
		Resource: "serviceaccount",
	}

	cmd := &cobra.Command{
		Use:     "serviceaccount NAME",
		Short:   "Get detailed information about a ServiceAccount",
		Long:    `Get detailed information about a ServiceAccount.`,
		Example: `potctl describe serviceaccount NAME`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			opt.Name = args[0]
			opt.Namespace, err = cmd.Flags().GetString("namespace")
			util.Check(err)

			exe, err := describe.NewExecutor(&opt)
			util.Check(err)

			err = exe.Execute()
			util.Check(err)
		},
	}
	cmd.Flags().StringVarP(&opt.Filename, "output-file", "o", "", "YAML output file")

	return cmd
}
