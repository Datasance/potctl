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

package deployagentconfig

import (
	// "errors"
	"fmt"
	"net/url"
	"strings"

	rsc "github.com/datasance/potctl/internal/resource"
	iutil "github.com/datasance/potctl/internal/util"
	clientutil "github.com/datasance/potctl/internal/util/client"
	"gopkg.in/yaml.v2"

	"github.com/datasance/iofog-go-sdk/v3/pkg/client"
	"github.com/datasance/potctl/internal/config"
	"github.com/datasance/potctl/internal/execute"
	"github.com/datasance/potctl/pkg/iofog"
	"github.com/datasance/potctl/pkg/iofog/install"
	"github.com/datasance/potctl/pkg/util"
)

type Options struct {
	Namespace string
	Yaml      []byte
	Name      string
	Tags      *[]string
}

type AgentConfigExecutor interface {
	execute.Executor
	GetAgentUUID() string
	SetHost(string)
	SetTags(*[]string)
	GetConfiguration() *rsc.AgentConfiguration
	GetNamespace() string
}

type RemoteExecutor struct {
	name        string
	uuid        string
	agentConfig *rsc.AgentConfiguration
	namespace   string
	tags        *[]string
}

func NewRemoteExecutor(name string, conf *rsc.AgentConfiguration, namespace string, tags *[]string) *RemoteExecutor {
	return &RemoteExecutor{
		name:        name,
		agentConfig: conf,
		namespace:   namespace,
		tags:        tags,
	}
}

func (exe *RemoteExecutor) GetNamespace() string {
	return exe.namespace
}

func (exe *RemoteExecutor) GetConfiguration() *rsc.AgentConfiguration {
	return exe.agentConfig
}

func (exe *RemoteExecutor) SetHost(host string) {
	exe.agentConfig.Host = &host
}

func (exe *RemoteExecutor) SetTags(tags *[]string) {
	// Merge tags
	if tags != nil {
		if exe.tags == nil {
			exe.tags = tags
		} else {
			newTagsSlice := append(*exe.tags, *tags...)
			exe.tags = &newTagsSlice
		}
	}
}

func (exe *RemoteExecutor) GetAgentUUID() string {
	return exe.uuid
}

func (exe *RemoteExecutor) GetName() string {
	return exe.name
}

func isOverridingSystemAgent(controllerHost, agentHost, agentName string) (err error) {
	// Generate controller endpoint
	controllerURL, err := url.Parse(controllerHost)
	if err != nil || controllerURL.Host == "" {
		controllerURL, err = url.Parse("//" + controllerHost) // Try to see if controllerEndpoint is an IP, in which case it needs to be pefixed by //
		if err != nil {
			return err
		}
	}
	agentURL, err := url.Parse(agentHost)
	if err != nil || agentURL.Host == "" {
		agentURL, err = url.Parse("//" + agentHost) // Try to see if controllerEndpoint is an IP, in which case it needs to be pefixed by //
		if err != nil {
			return err
		}
	}
	if agentURL.Hostname() == controllerURL.Hostname() && agentName != iofog.VanillaRemoteAgentName {
		return util.NewConflictError("Cannot deploy an agent on the same host than the Controller\n")
	}
	return nil
}

func (exe *RemoteExecutor) Execute() error {
	isSystem := iutil.IsSystemAgent(exe.agentConfig)

	// Check controller is reachable
	clt, err := clientutil.NewControllerClient(exe.namespace)
	if err != nil {
		return err
	}

	// Check we are not about to override Vanilla system agent
	ns, err := config.GetNamespace(exe.namespace)
	if err != nil {
		return err
	}
	controlPlane, err := ns.GetControlPlane()
	if err != nil || len(controlPlane.GetControllers()) == 0 {
		util.PrintError("You must deploy a Controller to a namespace before deploying any Agents")
		return err
	}

	// user := controlPlane.GetUser()

	// agents := ns.GetAgents()
	// numOfAgents := len(agents)

	// fmt.Println(": Current Number of Agents are ", numOfAgents)

	endpoint, err := controlPlane.GetEndpoint()
	if err != nil {
		fmt.Println("Error occurred while fetching endpoint from controlplane", err)
		return err
	}

	// baseURL, err := util.GetBaseURL(endpoint)
	// if err != nil {
	// 	fmt.Println("Error occurred while fetching base url from controlplane", err)
	// 	return err
	// }

	// ctrl, subscriptionKey, err := client.RefreshUserSubscriptionKey(client.Options{BaseURL: baseURL}, clt.GetRefreshToken(), user.Email, user.GetRawPassword())

	// if err != nil {
	// 	fmt.Println("Error occurred while fetching subscription key from controlplane: ", err)
	// 	return err
	// }

	// if ctrl == nil {
	// 	fmt.Println("Client came empty while fetching subscription key from controlplane")
	// 	return errors.New("Client came empty while fetching subscription key from controlplane")
	// }

	// if subscriptionKey != "" {
	// 	if user.SubscriptionKey != subscriptionKey {
	// 		fmt.Println("Subscription Key is updated from controlplane endpoints:", subscriptionKey)
	// 		user.SubscriptionKey = subscriptionKey
	// 		// controlPlane.UpdateUserSubscriptionKey(subscriptionKey)
	// 		config.UpdateSubscriptionKey(ns.Name, subscriptionKey)
	// 	}
	// }

	// expiryDate, agentSeats, err := util.GetEntitlementDatasance(user.SubscriptionKey, exe.namespace, user.Email)

	// if err != nil {
	// 	return err
	// }

	// isValid, reason := util.CheckExpiryDate(expiryDate)
	// if !isValid {
	// 	switch reason {
	// 	case "expired":
	// 		// This is the expected case for expired subscriptions.
	// 		return errors.New("subscription has expired")
	// 	case "not_found":
	// 		// Handle case when subscription is not found.
	// 		return errors.New("subscription not found")
	// 	case "engine_not_responding":
	// 		// Handle case when the subscription engine is not responding.
	// 		return errors.New("subscription engine is not responding")
	// 	default:
	// 		// Default case for unexpected errors, or if there's no specific reason.
	// 		return errors.New("unknown error with subscription expiry date check")
	// 	}
	// }

	// if !util.HasAvailableAgentSeats(numOfAgents, agentSeats) {
	// 	return errors.New("subscription limit reached: not enough agent seats available")
	// }

	if !isSystem || install.IsVerbose() {
		util.SpinStart(fmt.Sprintf("Deploying agent %s configuration", exe.GetName()))
	}

	host := ""
	if exe.agentConfig.Host != nil {
		host = *exe.agentConfig.Host
	}

	if err := isOverridingSystemAgent(endpoint, host, exe.name); err != nil {
		return err
	}

	// Get the Agent in question
	agent, err := clt.GetAgentByName(exe.name, isSystem)
	// TODO: replace this check with built-in IsNewNotFound() func from go-sdk
	if err != nil && !strings.Contains(err.Error(), "not find agent") {
		return err
	}
	ip := ""
	if agent != nil {
		ip = agent.IPAddressExternal
	}
	// Get all other non-system Agents
	agentList, err := clt.ListAgents(client.ListAgentsRequest{})
	if err != nil {
		return err
	}
	// Process needs to be done at execute time because agent might have been created during deploy
	if err := Process(exe.agentConfig, exe.name, ip, agentList.Agents); err != nil {
		return err
	}

	// Create if Agent does not exist
	if agent == nil {
		uuid, err := createAgentFromConfiguration(exe.agentConfig, exe.tags, exe.name, clt)
		exe.uuid = uuid
		return err
	}
	// Update existing Agent
	exe.uuid = agent.UUID
	return updateAgentConfiguration(exe.agentConfig, exe.tags, agent.UUID, clt)
}

func NewExecutor(opt Options) (exe execute.Executor, err error) {
	// Unmarshal file
	agentConfig := rsc.AgentConfiguration{}
	if err = yaml.UnmarshalStrict(opt.Yaml, &agentConfig); err != nil {
		err = util.NewUnmarshalError(err.Error())
		return
	}

	if agentConfig.Name == "" {
		agentConfig.Name = opt.Name
	}

	if err = Validate(&agentConfig); err != nil {
		return
	}

	return &RemoteExecutor{
		name:        opt.Name,
		agentConfig: &agentConfig,
		namespace:   opt.Namespace,
		tags:        opt.Tags,
	}, nil
}
