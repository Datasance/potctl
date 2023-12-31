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

package detachagent

import (
	"fmt"

	"github.com/datasance/potctl/pkg/iofog/install"
	"github.com/datasance/potctl/pkg/util"
)

func (exe executor) localDeprovision() error {
	containerClient, err := install.NewLocalContainerClient()
	if err != nil {
		util.PrintNotify(fmt.Sprintf("Could not deprovision local iofog-agent container. Error: %s\n", err.Error()))
	} else if _, err = containerClient.ExecuteCmd(install.GetLocalContainerName("agent", false), []string{
		"sudo",
		"iofog-agent",
		"deprovision",
	}); err != nil {
		util.PrintNotify(fmt.Sprintf("Could not deprovision local iofog-agent container. Error: %s\n", err.Error()))
	}
	return nil
}
