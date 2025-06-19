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

package deployremotecontrolplane

import (
	"fmt"
	// "strings"

	"github.com/datasance/iofog-go-sdk/v3/pkg/client"
	"github.com/datasance/potctl/internal/config"
	deployagent "github.com/datasance/potctl/internal/deploy/agent"
	deployagentconfig "github.com/datasance/potctl/internal/deploy/agentconfig"
	deployremotecontroller "github.com/datasance/potctl/internal/deploy/controller/remote"
	"github.com/datasance/potctl/internal/execute"
	rsc "github.com/datasance/potctl/internal/resource"
	iutil "github.com/datasance/potctl/internal/util"

	// clientutil "github.com/datasance/potctl/internal/util/client"
	"github.com/datasance/potctl/pkg/iofog"
	"github.com/datasance/potctl/pkg/iofog/install"
	"github.com/datasance/potctl/pkg/util"
)

type Options struct {
	Namespace string
	Yaml      []byte
	Name      string
}

type remoteControlPlaneExecutor struct {
	ctrlClient          *client.Client
	controllerExecutors []execute.Executor
	controlPlane        rsc.ControlPlane
	ns                  *rsc.Namespace
	name                string
}

func deploySystemAgent(namespace string, ctrl *rsc.RemoteController, systemAgent rsc.Package) (err error) {
	// Deploy system agent to host internal router
	install.Verbose("Deploying system agent for controller " + ctrl.Name)
	var deploymentType string
	if systemAgent.Container.Image != "" {
		deploymentType = "container"
	} else {
		deploymentType = "native"
	}
	agent := rsc.RemoteAgent{
		Name:    iofog.VanillaRemoteAgentName,
		Host:    ctrl.Host,
		SSH:     ctrl.SSH,
		Package: systemAgent,
	}
	// Configure agent to be system agent with default router
	RouterConfig := client.RouterConfig{
		RouterMode:      iutil.MakeStrPtr("interior"),
		MessagingPort:   iutil.MakeIntPtr(5671),
		EdgeRouterPort:  iutil.MakeIntPtr(45671),
		InterRouterPort: iutil.MakeIntPtr(55671),
	}

	upstreamRouters := []string{}

	deployAgentConfig := rsc.AgentConfiguration{
		Name:    iofog.VanillaRemoteAgentName,
		FogType: iutil.MakeStrPtr("auto"),
		AgentConfiguration: client.AgentConfiguration{
			IsSystem:        iutil.MakeBoolPtr(true),
			DeploymentType:  iutil.MakeStrPtr(deploymentType),
			Host:            &ctrl.Host,
			RouterConfig:    RouterConfig,
			UpstreamRouters: &upstreamRouters,
		},
	}

	// Get Agentconfig executor
	deployAgentConfigExecutor := deployagentconfig.NewRemoteExecutor(iofog.VanillaRemoteAgentName, &deployAgentConfig, namespace, nil)
	// If there already is a system fog, ignore error
	if err := deployAgentConfigExecutor.Execute(); err != nil {
		return err
	}
	agent.UUID = deployAgentConfigExecutor.GetAgentUUID()
	agentDeployExecutor, err := deployagent.NewRemoteExecutor(namespace, &agent, false)
	if err != nil {
		return err
	}
	return agentDeployExecutor.Execute()
}

func deployNonSystemAgent(namespace string, ctrl *rsc.RemoteController, systemAgent rsc.Package) (err error) {
	// Deploy system agent to host internal router
	install.Verbose("Deploying non-system agent for controller " + ctrl.Name)
	nonSystemAgentName := fmt.Sprintf("%s-controller", ctrl.Name)
	var deploymentType string
	if systemAgent.Container.Image != "" {
		deploymentType = "container"
	} else {
		deploymentType = "native"
	}
	agent := rsc.RemoteAgent{
		Name:    nonSystemAgentName,
		Host:    ctrl.Host,
		SSH:     ctrl.SSH,
		Package: systemAgent,
	}
	// Configure agent to be system agent with default router
	RouterConfig := client.RouterConfig{
		RouterMode:      iutil.MakeStrPtr("interior"),
		MessagingPort:   iutil.MakeIntPtr(5671),
		EdgeRouterPort:  iutil.MakeIntPtr(45671),
		InterRouterPort: iutil.MakeIntPtr(55671),
	}

	upstreamRouters := []string{"default-router"}

	deployAgentConfig := rsc.AgentConfiguration{
		Name:    nonSystemAgentName,
		FogType: iutil.MakeStrPtr("auto"),
		AgentConfiguration: client.AgentConfiguration{
			IsSystem:        iutil.MakeBoolPtr(false),
			DeploymentType:  iutil.MakeStrPtr(deploymentType),
			Host:            &ctrl.Host,
			RouterConfig:    RouterConfig,
			UpstreamRouters: &upstreamRouters,
		},
	}

	// Get Agentconfig executor
	deployAgentConfigExecutor := deployagentconfig.NewRemoteExecutor(iofog.VanillaRemoteAgentName, &deployAgentConfig, namespace, nil)
	// If there already is a system fog, ignore error
	if err := deployAgentConfigExecutor.Execute(); err != nil {
		return err
	}
	agent.UUID = deployAgentConfigExecutor.GetAgentUUID()
	agentDeployExecutor, err := deployagent.NewRemoteExecutor(namespace, &agent, false)
	if err != nil {
		return err
	}
	return agentDeployExecutor.Execute()
}

func tagControllerImage(ctrl *rsc.RemoteController, image string) (err error) {

	if image == "" {
		image = util.GetControllerImage()
	}

	// Connect
	ssh, err := util.NewSecureShellClient(ctrl.SSH.User, ctrl.Host, ctrl.SSH.KeyFile)
	if err != nil {
		return err
	}
	if err := ssh.Connect(); err != nil {
		return err
	}

	defer util.Log(ssh.Disconnect)

	cmds := []string{
		fmt.Sprintf(`echo "IOFOG_CONTROLLER_IMAGE=%s" | sudo tee -a "/etc/iofog/agent/iofog-agent.env" > /dev/null`, image),
		fmt.Sprintf("sudo service iofog-agent restart"),
	}

	// Execute commands
	for _, cmd := range cmds {
		_, err = ssh.Run(cmd)
		if err != nil {
			return
		}
	}

	return
}

func (exe remoteControlPlaneExecutor) postDeploy() (err error) {
	controllers := exe.controlPlane.GetControllers()
	// Deploy agents for each controller
	for idx, baseController := range controllers {
		controller, ok := baseController.(*rsc.RemoteController)
		if !ok {
			return util.NewInternalError("Could not convert Controller to Remote Controller")
		}
		remoteControlPlane, ok := exe.controlPlane.(*rsc.RemoteControlPlane)
		if !ok {
			return util.NewInternalError("Could not convert ControlPlane to Remote ControlPlane")
		}
		// First controller gets system agent, others get non-system agents
		if idx == 0 {
			if err := deploySystemAgent(exe.ns.Name, controller, remoteControlPlane.SystemAgent); err != nil {
				return fmt.Errorf("failed to deploy system agent for first controller: %v", err)
			}
		} else {
			if err := deployNonSystemAgent(exe.ns.Name, controller, remoteControlPlane.SystemAgent); err != nil {
				return fmt.Errorf("failed to deploy non-system agent for controller %d: %v", idx, err)
			}
		}

		// Tag controller image for all controllers
		if err := tagControllerImage(controller, remoteControlPlane.Package.Container.Image); err != nil {
			return fmt.Errorf("failed to tag controller image for controller %d: %v", idx, err)
		}
	}
	return nil
}

func (exe remoteControlPlaneExecutor) Execute() (err error) {
	util.SpinStart(fmt.Sprintf("Deploying controlplane %s", exe.GetName()))
	if err := runExecutors(exe.controllerExecutors); err != nil {
		return err
	}

	// Make sure Controller API is ready
	endpoint, err := exe.controlPlane.GetEndpoint()
	if err != nil {
		return
	}
	if err := install.WaitForControllerAPI(endpoint); err != nil {
		return err
	}
	// // Create new user
	// baseURL, err := util.GetBaseURL(endpoint)
	// if err != nil {
	// 	return err
	// }
	// exe.ctrlClient = client.New(client.Options{BaseURL: baseURL})
	// user := client.User(exe.controlPlane.GetUser())
	// user.Password = exe.controlPlane.GetUser().GetRawPassword()
	// if err = exe.ctrlClient.CreateUser(user); err != nil {
	// 	// If not error about account existing, fail
	// 	if !strings.Contains(err.Error(), "already an account associated") {
	// 		return err
	// 	}
	// 	// Try to log in
	// 	if err := exe.ctrlClient.Login(client.LoginRequest{
	// 		Email:    user.Email,
	// 		Password: user.Password,
	// 	}); err != nil {
	// 		return err
	// 	}
	// }
	// Update config
	exe.ns.SetControlPlane(exe.controlPlane)
	if err := config.Flush(); err != nil {
		return err
	}

	// Post deploy steps
	return exe.postDeploy()
}

func (exe remoteControlPlaneExecutor) GetName() string {
	return exe.name
}

func newControlPlaneExecutor(executors []execute.Executor, namespace *rsc.Namespace, name string, controlPlane rsc.ControlPlane) execute.Executor {
	return remoteControlPlaneExecutor{
		controllerExecutors: executors,
		ns:                  namespace,
		controlPlane:        controlPlane,
		name:                name,
	}
}

// Validates database configuration for multi-controller setup
func validateMultiControllerDatabase(controlPlane *rsc.RemoteControlPlane) error {
	if len(controlPlane.Controllers) > 1 {
		db := controlPlane.Database
		if db.Provider == "" || db.Host == "" || db.DatabaseName == "" ||
			db.Password == "" || db.Port == 0 || db.User == "" {
			return util.NewInputError("When deploying multiple controllers, you must specify an external database configuration with all required fields (host, user, password, provider, databaseName, port)")
		}
	}
	return nil
}

// Validates HTTPS configuration for a single controller
func validateControllerHTTPS(controller *rsc.RemoteController) error {
	if controller.Https != nil && controller.Https.Enabled != nil && *controller.Https.Enabled {
		// HTTPS is enabled, validate required fields
		if controller.Https.TLSCert == "" || controller.Https.TLSKey == "" {
			return util.NewInputError("When HTTPS is enabled, you must provide TLS certificate and key")
		}
	}
	return nil
}

// Validates CA configuration for a controller
func validateControllerRouterCA(controller *rsc.RemoteController) error {
	if controller.SiteCA != nil {
		if controller.SiteCA.TLSCert == "" || controller.SiteCA.TLSKey == "" {
			return util.NewInputError("When SiteCA is configured, you must provide both TLS certificate and key")
		}
	}
	if controller.LocalCA != nil {
		if controller.LocalCA.TLSCert == "" || controller.LocalCA.TLSKey == "" {
			return util.NewInputError("When LocalCA is configured, you must provide both TLS certificate and key")
		}
	}
	return nil
}

// Validates HTTPS configuration across all controllers
func validateMultiControllerHTTPS(controlPlane *rsc.RemoteControlPlane) error {
	controllers := controlPlane.Controllers
	if len(controllers) <= 1 {
		return nil
	}

	// Check first controller's HTTPS config
	firstController := controllers[0]
	if firstController.Https != nil && firstController.Https.Enabled != nil && *firstController.Https.Enabled {
		// First controller has HTTPS enabled, validate all controllers
		for idx, controller := range controllers {
			if err := validateControllerHTTPS(&controller); err != nil {
				return fmt.Errorf("controller %d (%s): %v", idx, controller.Name, err)
			}
		}
	}
	return nil
}

// Validates CA configuration across all controllers
func validateMultiControllerRouterCA(controlPlane *rsc.RemoteControlPlane) error {
	controllers := controlPlane.Controllers
	if len(controllers) <= 1 {
		return nil
	}

	// Only first controller should have CA configuration
	firstController := controllers[0]
	if firstController.SiteCA != nil || firstController.LocalCA != nil {
		// Validate first controller's CA config
		if err := validateControllerRouterCA(&firstController); err != nil {
			return fmt.Errorf("first controller (%s): %v", firstController.Name, err)
		}

		// Check that other controllers don't have CA config
		for idx, controller := range controllers[1:] {
			if controller.SiteCA != nil || controller.LocalCA != nil {
				return fmt.Errorf("controller %d (%s): CA configuration should only be specified for the first controller", idx+1, controller.Name)
			}
		}
	}
	return nil
}

// Main validation function that orchestrates all validations
func validateMultiControllerConfig(controlPlane *rsc.RemoteControlPlane) error {
	// Validate database configuration
	if err := validateMultiControllerDatabase(controlPlane); err != nil {
		return err
	}

	// Validate HTTPS configuration
	if err := validateMultiControllerHTTPS(controlPlane); err != nil {
		return err
	}

	// Validate CA configuration
	if err := validateMultiControllerRouterCA(controlPlane); err != nil {
		return err
	}

	return nil
}

func NewExecutor(opt Options) (exe execute.Executor, err error) {
	// Check the namespace exists
	ns, err := config.GetNamespace(opt.Namespace)
	if err != nil {
		return
	}

	// Read the input file
	controlPlane, err := rsc.UnmarshallRemoteControlPlane(opt.Yaml)
	if err != nil {
		return
	}

	// Validate control plane for multiple controllers
	if err := validateMultiControllerConfig(&controlPlane); err != nil {
		return nil, err
	}

	// Create exe Controllers
	controllers := controlPlane.GetControllers()
	controllerExecutors := make([]execute.Executor, len(controllers))
	for idx := range controllers {
		controller, ok := controllers[idx].(*rsc.RemoteController)
		if !ok {
			return nil, util.NewError("Could not convert Controller to Remote Controller")
		}
		exe, err := deployremotecontroller.NewExecutorWithoutParsing(opt.Namespace, &controlPlane, controller)
		if err != nil {
			return nil, err
		}
		controllerExecutors[idx] = exe
	}

	return newControlPlaneExecutor(controllerExecutors, ns, opt.Name, &controlPlane), nil
}

func runExecutors(executors []execute.Executor) error {
	if errs, _ := execute.ForParallel(executors); len(errs) > 0 {
		return execute.CoalesceErrors(errs)
	}
	return nil
}
