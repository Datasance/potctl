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

package client

import (
	"fmt"

	"github.com/datasance/iofog-go-sdk/v3/pkg/client"
	"github.com/datasance/potctl/internal/config"
	rsc "github.com/datasance/potctl/internal/resource"
	"github.com/datasance/potctl/pkg/iofog"
	"github.com/datasance/potctl/pkg/util"
)

// clientCacheRoutine handles concurrent requests for a cached Controller client
func clientCacheRoutine() {
	for {
		request := <-pkg.clientCacheRequestChan
		// Invalidate cache
		if request.namespace == "" {
			pkg.clientCache = make(map[string]*client.Client)
			continue
		}
		result := &clientCacheResult{}
		// From cache
		if cachedClient, exists := pkg.clientCache[request.namespace]; exists {
			result.client = cachedClient
			request.resultChan <- result
			continue
		}
		// Create new client
		ioClient, err := newControllerClient(request.namespace)
		// Failure
		if err != nil {
			result.err = err
			request.resultChan <- result
			continue
		}
		// Save to cache and return new client
		pkg.clientCache[request.namespace] = ioClient
		result.client = ioClient
		request.resultChan <- result
	}
}

// agentCacheRoutine handles concurrent requests for a cached list of Agents
func agentCacheRoutine() {
	for {
		request := <-pkg.agentCacheRequestChan
		if request.namespace == "" {
			// Invalidate cache
			pkg.agentCache = make(map[string][]client.AgentInfo)
			continue
		}
		result := &agentCacheResult{}
		// From cache
		if cachedAgents, exist := pkg.agentCache[request.namespace]; exist {
			result.agents = cachedAgents
			request.resultChan <- result
			continue
		}
		// Client to get agents
		ioClient, err := NewControllerClient(request.namespace)
		if err != nil {
			result.err = err
			request.resultChan <- result
			continue
		}
		// Get agents
		agents, err := getBackendAgents(request.namespace, ioClient)
		if err != nil {
			result.err = err
			request.resultChan <- result
			continue
		}
		// Save to cache and return new agents
		pkg.agentCache[request.namespace] = agents
		result.agents = agents
		request.resultChan <- result
	}
}

func agentSyncRoutine() {
	complete := false
	for {
		request := <-pkg.agentSyncRequestChan
		if complete {
			request.resultChan <- nil
			continue
		}
		if err := syncAgentInfo(request.namespace); err != nil {
			request.resultChan <- err
			continue
		}
		complete = true
		request.resultChan <- nil
	}
}

func syncAgentInfo(namespace string) error {
	// Get local cache Agents
	ns, err := config.GetNamespace(namespace)
	if err != nil {
		return err
	}
	// Check the Control Plane type
	controlPlane, err := ns.GetControlPlane()
	if err != nil {
		return err
	}
	if _, ok := controlPlane.(*rsc.LocalControlPlane); ok {
		// Do not update local Agents
		return nil
	}
	// Generate map of config Agents
	agentsMap := make(map[string]*rsc.RemoteAgent)
	var localAgent *rsc.LocalAgent
	for _, baseAgent := range ns.GetAgents() {
		if v, ok := baseAgent.(*rsc.LocalAgent); ok {
			localAgent = v
		} else {
			agentsMap[baseAgent.GetName()] = baseAgent.(*rsc.RemoteAgent)
		}
	}

	// Get backend Agents
	backendAgents, err := GetBackendAgents(namespace)
	if err != nil {
		return err
	}

	// Generate cache types
	agents := make([]rsc.RemoteAgent, len(backendAgents))
	for idx := range backendAgents {
		backendAgent := &backendAgents[idx]
		if localAgent != nil && backendAgent.Name == localAgent.Name {
			localAgent.UUID = backendAgent.UUID
			continue
		}

		agent := rsc.RemoteAgent{
			Name: backendAgent.Name,
			UUID: backendAgent.UUID,
			Host: backendAgent.Host,
		}
		// Update additional info if local cache contains it
		if cachedAgent, exists := agentsMap[backendAgent.Name]; exists {
			agent.Created = cachedAgent.GetCreatedTime()
			agent.SSH = cachedAgent.SSH
		}

		agents[idx] = agent
	}

	// Overwrite the Agents
	ns.DeleteAgents()
	for idx := range agents {
		if err := ns.AddAgent(&agents[idx]); err != nil {
			return err
		}
	}

	if localAgent != nil {
		if err := ns.AddAgent(localAgent); err != nil {
			return err
		}
	}

	return config.Flush()
}

func newControllerClient(namespace string) (*client.Client, error) {

	// Try to get the client from the cache first
	if cachedClient, found := pkg.clientCache[namespace]; found {

		// If a cached client exists, use SessionLogin to refresh the session
		ns, err := config.GetNamespace(namespace)
		if err != nil {
			return nil, err
		}

		controlPlane, err := ns.GetControlPlane()
		if err != nil {
			return nil, err
		}

		endpoint, err := controlPlane.GetEndpoint()
		if err != nil {
			return nil, err
		}

		user := controlPlane.GetUser()

		// Get base URL
		baseURL, err := util.GetBaseURL(endpoint)
		if err != nil {
			return nil, err
		}

		// Use the refresh token from the cached client
		refreshToken := cachedClient.GetRefreshToken()
		user.AccessToken = cachedClient.GetAccessToken()
		user.RefreshToken = cachedClient.GetRefreshToken()
		// controlPlane.UpdateUserTokens(user.AccessToken, user.RefreshToken)
		config.UpdateUser(namespace, user.AccessToken, user.RefreshToken)

		// Use SessionLogin to attempt to refresh the session
		refreshedClient, err := client.SessionLogin(client.Options{BaseURL: baseURL}, refreshToken, user.Email, user.GetRawPassword())
		if err != nil {
			fmt.Println("Error: Failed to refresh session:", err)
			return nil, fmt.Errorf("failed to refresh session: %v", err)
		}

		// Update the cached client with the refreshed session
		pkg.clientCache[namespace] = refreshedClient
		config.Flush()
		return refreshedClient, nil
	}

	// If no cached client, proceed with NewAndLogin to create a new client
	ns, err := config.GetNamespace(namespace)
	if err != nil {
		return nil, err
	}
	controlPlane, err := ns.GetControlPlane()
	if err != nil {
		return nil, err
	}
	endpoint, err := controlPlane.GetEndpoint()
	if err != nil {
		return nil, err
	}

	user := controlPlane.GetUser()
	baseURL, err := util.GetBaseURL(endpoint)
	if err != nil {
		return nil, err
	}

	// Create a new client and login
	newClient, err := client.SessionLogin(client.Options{BaseURL: baseURL}, user.RefreshToken, user.Email, user.GetRawPassword())
	if err != nil {
		return nil, err
	}

	user.AccessToken = newClient.GetAccessToken()
	user.RefreshToken = newClient.GetRefreshToken()
	// controlPlane.UpdateUserTokens(user.AccessToken, user.RefreshToken)
	config.UpdateUser(namespace, user.AccessToken, user.RefreshToken)

	// Flush the config and handle errors
	if err := config.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush config: %v", err)
	}

	return newClient, nil
}

func getBackendAgents(namespace string, ioClient *client.Client) ([]client.AgentInfo, error) {
	agentList, err := ioClient.ListAgents(client.ListAgentsRequest{})
	if err != nil {
		return nil, err
	}
	pkg.agentCache[namespace] = agentList.Agents
	return agentList.Agents, nil
}

func getAgentNameFromUUID(agentMapByUUID map[string]client.AgentInfo, uuid string) (name string) {
	if uuid == iofog.VanillaRemoteAgentName {
		return uuid
	}
	if uuid == iofog.VanillaRouterAgentName {
		return uuid
	}
	agent, found := agentMapByUUID[uuid]
	if !found {
		util.PrintNotify(fmt.Sprintf("Could not find Router: %s\n", uuid))
		name = "UNKNOWN ROUTER: " + uuid
	} else {
		name = agent.Name
	}
	return
}
