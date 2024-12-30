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

package install

import (
	"github.com/datasance/iofog-go-sdk/v3/pkg/client"
	"github.com/datasance/potctl/pkg/util"
)

type Agent interface {
	Bootstrap() error
	getProvisionKey(string, IofogUser) (string, string, error)
}

// defaultAgent implements commong behavior
type defaultAgent struct {
	name string
	uuid string
}

func (agent *defaultAgent) getProvisionKey(controllerEndpoint string, user IofogUser) (key string, err error) {
	// Connect to controller
	baseURL, err := util.GetBaseURL(controllerEndpoint)
	if err != nil {
		return
	}
	// Log in
	ctrl, err := client.SessionLogin(client.Options{BaseURL: baseURL}, user.RefreshToken, user.Email, user.Password)
	if err != nil {
		return
	}

	Verbose("Accessing Controller to generate Provisioning Key")
	// loginRequest := client.LoginRequest{
	// 	Email:    user.Email,
	// 	Password: user.Password,
	// }
	// if err = ctrl.Login(loginRequest); err != nil {
	// 	return
	// }

	// System agents have uuid passed through, normal agents dont
	if agent.uuid == "" {
		var agentInfo *client.AgentInfo
		agentInfo, err = ctrl.GetAgentByName(agent.name, false)
		if err != nil {
			return
		}
		agent.uuid = agentInfo.UUID
	}

	// Get provisioning key
	provisionResponse, err := ctrl.GetAgentProvisionKey(agent.uuid)
	if err != nil {
		return
	}
	key = provisionResponse.Key
	return key, err
}
