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

package deployagent

import (
	"fmt"

	"github.com/datasance/potctl/internal/config"
	rsc "github.com/datasance/potctl/internal/resource"
	// clientutil "github.com/datasance/potctl/internal/util/client"
	"github.com/datasance/potctl/pkg/iofog"
	"github.com/datasance/potctl/pkg/iofog/install"
	"github.com/datasance/potctl/pkg/util"
)

type remoteExecutor struct {
	namespace string
	agent     *rsc.RemoteAgent
}

func newRemoteExecutor(namespace string, agent *rsc.RemoteAgent) *remoteExecutor {
	exe := &remoteExecutor{}
	exe.namespace = namespace
	exe.agent = agent

	return exe
}

func (exe *remoteExecutor) GetName() string {
	return exe.agent.Name
}

func (exe *remoteExecutor) ProvisionAgent() (string, error) {
	// Get agent
	// agent, err := install.NewRemoteAgent(exe.agent.SSH.User,
	// 	exe.agent.Host,
	// 	exe.agent.SSH.Port,
	// 	exe.agent.SSH.KeyFile,
	// 	exe.agent.Name,
	// 	exe.agent.UUID)
	var agent *install.RemoteAgent
	var err error

	if exe.agent.Package.Container.Image != "" {
		// Use NewRemoteContainerAgent
		agent, err = install.NewRemoteContainerAgent(
			exe.agent.SSH.User,
			exe.agent.Host,
			exe.agent.SSH.Port,
			exe.agent.SSH.KeyFile,
			exe.agent.Name,
			exe.agent.UUID,
		)
	} else {
		// Use NewRemoteAgent
		agent, err = install.NewRemoteAgent(
			exe.agent.SSH.User,
			exe.agent.Host,
			exe.agent.SSH.Port,
			exe.agent.SSH.KeyFile,
			exe.agent.Name,
			exe.agent.UUID,
		)
	}

	if err != nil {
		return "", err
	}

	ns, err := config.GetNamespace(exe.namespace)
	if err != nil {
		return "", err
	}
	controlPlane, err := ns.GetControlPlane()
	if err != nil {
		return "", err
	}
	// Try Agent-specific endpoint first
	controllerEndpoint := exe.agent.GetControllerEndpoint()
	if controllerEndpoint == "" {
		controllerEndpoint, err = controlPlane.GetEndpoint()
		if err != nil {
			return "", util.NewError("Failed to retrieve Controller endpoint!")
		}
	}

	// Configure the agent with Controller details
	user := install.IofogUser(controlPlane.GetUser())
	user.Password = controlPlane.GetUser().GetRawPassword()
	// Get Config before provision and set iofog-agent config
	agentConfig := exe.agent.GetConfig()

	// Check if agentConfig is empty
	if agentConfig == nil || (agentConfig.Name == "" && agentConfig.FogType == nil) {
		util.PrintNotify(fmt.Sprintf("Skipping initial agent configuration for %s as agent config parameters are empty. Default config parameters will be used.", exe.agent.Name))

	} else {
		err = agent.SetInitialConfig(
			agentConfig.Name,
			agentConfig.Location,
			agentConfig.Latitude,
			agentConfig.Longitude,
			agentConfig.Description,
			*agentConfig.FogType,           // Dereference after checking nil
			agentConfig.AgentConfiguration, // Pass the embedded client.AgentConfiguration
		)
		if err != nil {
			return "", err
		}
	}
	// err = agent.SetInitialConfig(
	// 	agentConfig.Name,
	// 	agentConfig.Location,
	// 	agentConfig.Latitude,
	// 	agentConfig.Longitude,
	// 	agentConfig.Description,
	// 	*agentConfig.FogType,
	// 	agentConfig.AgentConfiguration, // Pass the embedded client.AgentConfiguration
	// )
	// if err != nil {
	// 	return "", err
	// }
	return agent.Configure(controllerEndpoint, user)
}

// Deploy iofog-agent stack on an agent host
func (exe *remoteExecutor) Execute() (err error) {
	// Get Control Plane
	ns, err := config.GetNamespace(exe.namespace)
	if err != nil {
		return err
	}
	controlPlane, err := ns.GetControlPlane()
	if err != nil || len(controlPlane.GetControllers()) == 0 {
		util.PrintError("You must deploy a Controller to a namespace before deploying any Agents")
		return
	}

	// Connect to agent via SSH
	// agent, err := install.NewRemoteAgent(exe.agent.SSH.User,
	// 	exe.agent.Host,
	// 	exe.agent.SSH.Port,
	// 	exe.agent.SSH.KeyFile,
	// 	exe.agent.Name,
	// 	exe.agent.UUID)

	var agent *install.RemoteAgent

	if exe.agent.Package.Container.Image != "" {
		// Use NewRemoteContainerAgent
		agent, err = install.NewRemoteContainerAgent(
			exe.agent.SSH.User,
			exe.agent.Host,
			exe.agent.SSH.Port,
			exe.agent.SSH.KeyFile,
			exe.agent.Name,
			exe.agent.UUID,
		)
	} else {
		// Use NewRemoteAgent
		agent, err = install.NewRemoteAgent(
			exe.agent.SSH.User,
			exe.agent.Host,
			exe.agent.SSH.Port,
			exe.agent.SSH.KeyFile,
			exe.agent.Name,
			exe.agent.UUID,
		)
	}
	if err != nil {
		return err
	}

	// Set custom scripts
	if exe.agent.Scripts != nil {
		if err := agent.CustomizeProcedures(
			exe.agent.Scripts.Directory,
			&exe.agent.Scripts.AgentProcedures); err != nil {
			return err
		}
	}
	if exe.agent.Package.Container.Image != "" {
		// Set Image
		agent.SetContainerImage(exe.agent.Package.Container.Image)

	} else {
		// Set version
		agent.SetVersion(exe.agent.Package.Version)
	}
	// // Set version
	// agent.SetVersion(exe.agent.Package.Version)
	// agent.SetRepository(exe.agent.Package.Repo, exe.agent.Package.Token)

	// Try the deploy
	err = agent.Bootstrap()
	if err != nil {
		return
	}

	uuid, err := exe.ProvisionAgent()
	if err != nil {
		return err
	}

	// Return the Agent through pointer
	exe.agent.UUID = uuid
	exe.agent.Created = util.NowUTC()
	return nil
}

func ValidateRemoteAgent(agent *rsc.RemoteAgent) error {
	if err := util.IsLowerAlphanumeric("Agent", agent.Name); err != nil {
		return err
	}
	if agent.Name == iofog.VanillaRouterAgentName || agent.Name == iofog.VanillaRemoteAgentName {
		return util.NewInputError(fmt.Sprintf("%s is a reserved name and cannot be used for an Agent", iofog.VanillaRouterAgentName))
	}
	if (agent.Host != "localhost" && agent.Host != "127.0.0.1") && (agent.Host == "" || agent.SSH.User == "" || agent.SSH.KeyFile == "") {
		return util.NewInputError("For Agents you must specify non-empty values for host, user, and keyfile")
	}
	return nil
}
