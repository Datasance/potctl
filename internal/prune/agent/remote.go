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

package pruneagent

import (
	"fmt"
	"strings"

	rsc "github.com/datasance/potctl/v1/internal/resource"
	clientutil "github.com/datasance/potctl/v1/internal/util/client"
	"github.com/datasance/potctl/v1/pkg/iofog/install"
	"github.com/datasance/potctl/v1/pkg/util"
)

func (exe executor) remoteAgentPrune(agent rsc.Agent) error {
	ctrl, err := clientutil.NewControllerClient(exe.namespace)
	if err != nil {
		return err
	}
	// If controller exists, prune the agent
	// Perform Docker pruning of Agent through Controller
	if err = ctrl.PruneAgent(agent.GetUUID()); err != nil {
		if !strings.Contains(err.Error(), "NotFoundError") {
			return err
		}
	}
	return nil
}

func (exe executor) remoteDetachedAgentPrune(agent *rsc.RemoteAgent) error {
	if err := agent.ValidateSSH(); err != nil {
		return err
	}
	sshAgent, err := install.NewRemoteAgent(agent.SSH.User, agent.Host, agent.SSH.Port, agent.SSH.KeyFile, agent.Name, agent.UUID)
	if err != nil {
		return err
	}
	if err := sshAgent.Prune(); err != nil {
		return util.NewInternalError(fmt.Sprintf("Failed to Prune Iofog resource %s. %s", agent.Name, err.Error()))
	}
	return nil
}
