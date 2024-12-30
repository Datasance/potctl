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

package deploylocalcontrolplane

import (
	"fmt"
	// "strings"

	"github.com/datasance/iofog-go-sdk/v3/pkg/client"
	"github.com/datasance/potctl/internal/config"
	// deployagentconfig "github.com/datasance/potctl/internal/deploy/agentconfig"
	deploylocalcontroller "github.com/datasance/potctl/internal/deploy/controller/local"
	"github.com/datasance/potctl/internal/execute"
	rsc "github.com/datasance/potctl/internal/resource"
	iutil "github.com/datasance/potctl/internal/util"
	clientutil "github.com/datasance/potctl/internal/util/client"
	// "github.com/datasance/potctl/pkg/iofog"
	"github.com/datasance/potctl/pkg/iofog/install"
	"github.com/datasance/potctl/pkg/util"
)

type Options struct {
	Namespace string
	Yaml      []byte
	Name      string
}
type localControlPlaneExecutor struct {
	ctrlClient          *client.Client
	controllerExecutors []execute.Executor
	controlPlane        *rsc.LocalControlPlane
	namespace           string
	name                string
}

func createDefaultRouter(clt *client.Client) (err error) {
	routerConfig := client.Router{
		Host: "localhost",
		RouterConfig: client.RouterConfig{
			RouterMode:      iutil.MakeStrPtr("interior"),
			MessagingPort:   iutil.MakeIntPtr(5672),
			EdgeRouterPort:  iutil.MakeIntPtr(56721),
			InterRouterPort: iutil.MakeIntPtr(56722),
		},
	}

	return clt.PutDefaultRouter(routerConfig)
}

func (exe localControlPlaneExecutor) postDeploy() (err error) {
	// Check controller is reachable
	clt, err := clientutil.NewControllerClient(exe.namespace)
	if err != nil {
		return err
	}

	if err := createDefaultRouter(clt); err != nil {
		return err
	}
	return nil
}

func (exe localControlPlaneExecutor) Execute() (err error) {
	util.SpinStart(fmt.Sprintf("Deploying controlplane %s", exe.GetName()))
	if err := runExecutors(exe.controllerExecutors); err != nil {
		return err
	}

	// Make sure Controller API is ready
	controller, err := exe.controlPlane.GetController("")
	if err != nil {
		return err
	}
	endpoint := controller.GetEndpoint()

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
	ns, err := config.GetNamespace(exe.namespace)
	if err != nil {
		return err
	}
	ns.SetControlPlane(exe.controlPlane)
	if err := config.Flush(); err != nil {
		return err
	}
	// Post deploy steps
	return exe.postDeploy()
}

func (exe localControlPlaneExecutor) GetName() string {
	return exe.name
}

func newControlPlaneExecutor(executors []execute.Executor, namespace, name string, controlPlane *rsc.LocalControlPlane) execute.Executor {
	return localControlPlaneExecutor{
		controllerExecutors: executors,
		namespace:           namespace,
		controlPlane:        controlPlane,
		name:                name,
	}
}

func NewExecutor(opt Options) (exe execute.Executor, err error) {
	// Check the namespace exists
	_, err = config.GetNamespace(opt.Namespace)
	if err != nil {
		return
	}

	// Read the input file
	controlPlane, err := rsc.UnmarshallLocalControlPlane(opt.Yaml)
	if err != nil {
		return
	}

	// Create exe Controllers
	controllers := controlPlane.GetControllers()
	controllerExecutors := make([]execute.Executor, len(controllers))
	for idx := range controllers {
		controller, ok := controllers[idx].(*rsc.LocalController)
		if !ok {
			return nil, util.NewError("Could not convert Controller to Local Controller")
		}
		exe, err := deploylocalcontroller.NewExecutorWithoutParsing(opt.Namespace, &controlPlane, controller)
		if err != nil {
			return nil, err
		}
		controllerExecutors[idx] = exe
	}

	return newControlPlaneExecutor(controllerExecutors, opt.Namespace, opt.Name, &controlPlane), nil
}

func runExecutors(executors []execute.Executor) error {
	if errs, _ := execute.ForParallel(executors); len(errs) > 0 {
		return execute.CoalesceErrors(errs)
	}
	return nil
}
