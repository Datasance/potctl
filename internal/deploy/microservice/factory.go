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

package deploymicroservice

import (
	"fmt"

	apps "github.com/datasance/iofog-go-sdk/v3/pkg/apps"
	"github.com/datasance/potctl/internal/config"
	"github.com/datasance/potctl/internal/execute"
	clientutil "github.com/datasance/potctl/internal/util/client"
	"github.com/datasance/potctl/pkg/util"
	"gopkg.in/yaml.v2"
)

type Options struct {
	Namespace string
	Yaml      []byte
	Name      string
}

type remoteExecutor struct {
	namespace    string
	microservice interface{}
	name         string
}

func (exe *remoteExecutor) GetName() string {
	return exe.name
}

func (exe *remoteExecutor) Execute() error {
	util.SpinStart(fmt.Sprintf("Deploying microservice %s", exe.GetName()))
	ns, err := config.GetNamespace(exe.namespace)
	if err != nil {
		return err
	}
	controlPlane, err := ns.GetControlPlane()
	if err != nil {
		return err
	}

	// Check Controller exists
	if len(controlPlane.GetControllers()) == 0 {
		return util.NewInputError("This namespace does not have a Controller. You must first deploy a Controller before deploying Applications")
	}
	endpoint, err := controlPlane.GetEndpoint()
	if err != nil {
		return err
	}

	clt, err := clientutil.NewControllerClient(exe.namespace)
	if err != nil {
		return err
	}

	controller := apps.IofogController{
		Endpoint:     endpoint,
		Email:        controlPlane.GetUser().Email,
		Password:     controlPlane.GetUser().Password,
		Token:        clt.GetAccessToken(),
		RefreshToken: clt.GetRefreshToken(),
	}

	appName, msvcName, err := clientutil.ParseFQName(exe.name, "Microservice")
	if err != nil {
		return err
	}

	return apps.DeployMicroservice(controller, exe.microservice, appName, msvcName)
}

func NewExecutor(opt Options) (exe execute.Executor, err error) {
	// Check the namespace exists
	if _, err = config.GetNamespace(opt.Namespace); err != nil {
		return exe, err
	}
	// Unmarshal file
	var microservice interface{}
	if err = yaml.UnmarshalStrict(opt.Yaml, &microservice); err != nil {
		err = util.NewUnmarshalError(err.Error())
		return
	}

	return &remoteExecutor{
		namespace:    opt.Namespace,
		microservice: &microservice,
		name:         opt.Name,
	}, nil
}
