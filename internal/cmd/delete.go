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
	"errors"

	"github.com/datasance/potctl/internal/delete"
	"github.com/datasance/potctl/pkg/util"
	"github.com/spf13/cobra"
)

func newDeleteCommand() *cobra.Command {
	// Instantiate options
	opt := &delete.Options{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an existing ioFog resource",
		Long:  `Delete an existing ioFog resource.`,
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			opt.Namespace, err = cmd.Flags().GetString("namespace")
			util.Check(err)

			// Check file
			if opt.InputFile == "" {
				util.Check(errors.New("provided empty value for input file via the -f flag"))
			}

			// Execute command
			err = delete.Execute(opt)
			util.Check(err)

			util.PrintSuccess("Successfully deleted resources from namespace " + opt.Namespace)
		},
	}

	// Add subcommands
	cmd.AddCommand(
		newDeleteNamespaceCommand(),
		newDeleteControllerCommand(),
		newDeleteAgentCommand(),
		newDeleteAllCommand(),
		newDeleteApplicationCommand(),
		newDeleteApplicationTemplateCommand(),
		newDeleteCatalogItemCommand(),
		newDeleteRegistryCommand(),
		newDeleteMicroserviceCommand(),
		newDeleteVolumeCommand(),
		newDeleteRouteCommand(),
		newDeleteEdgeResourceCommand(),
	)

	// Register flags
	cmd.Flags().StringVarP(&opt.InputFile, "file", "f", "", pkg.flagDescYaml)

	return cmd
}
