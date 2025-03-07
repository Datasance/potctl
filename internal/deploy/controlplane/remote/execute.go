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
	clientutil "github.com/datasance/potctl/internal/util/client"
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
	install.Verbose("Deploying system agent")
	agent := rsc.RemoteAgent{
		Name:    iofog.VanillaRemoteAgentName,
		Host:    ctrl.Host,
		SSH:     ctrl.SSH,
		Package: systemAgent,
	}
	// Configure agent to be system agent with default router
	RouterConfig := client.RouterConfig{
		RouterMode:      iutil.MakeStrPtr("interior"),
		MessagingPort:   iutil.MakeIntPtr(5672),
		EdgeRouterPort:  iutil.MakeIntPtr(56721),
		InterRouterPort: iutil.MakeIntPtr(56722),
	}

	upstreamRouters := []string{}

	deployAgentConfig := rsc.AgentConfiguration{
		Name:    iofog.VanillaRemoteAgentName,
		FogType: iutil.MakeStrPtr("auto"),
		AgentConfiguration: client.AgentConfiguration{
			IsSystem:        iutil.MakeBoolPtr(false),
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

func createDefaultRouter(namespace string, ctrl *rsc.RemoteController) (err error) {
	// Check controller is reachable
	clt, err := clientutil.NewControllerClient(namespace)
	if err != nil {
		return err
	}

	routerConfig := client.Router{
		Host: ctrl.Host,
		RouterConfig: client.RouterConfig{
			RouterMode:      iutil.MakeStrPtr("interior"),
			MessagingPort:   iutil.MakeIntPtr(5672),
			EdgeRouterPort:  iutil.MakeIntPtr(56721),
			InterRouterPort: iutil.MakeIntPtr(56722),
		},
	}

	if err := clt.PutDefaultRouter(routerConfig); err != nil {
		return err
	}
	return nil
}

func (exe remoteControlPlaneExecutor) postDeploy() (err error) {
	// Look for a Vanilla controller
	controllers := exe.controlPlane.GetControllers()
	for _, baseController := range controllers {
		controller, ok := baseController.(*rsc.RemoteController)
		if !ok {
			return util.NewInternalError("Could not convert Controller to Remote Controller")
		}
		remoteControlPlane, ok := exe.controlPlane.(*rsc.RemoteControlPlane)
		if !ok {
			return util.NewInternalError("Could not convert ControlPlane to Remote ControlPlane")
		}
		if err := createDefaultRouter(exe.ns.Name, controller); err != nil {
			return err
		}
		if err := deploySystemAgent(exe.ns.Name, controller, remoteControlPlane.SystemAgent); err != nil {
			return err
		}
		if err := tagControllerImage(controller, remoteControlPlane.Package.Container.Image); err != nil {
			return err
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
