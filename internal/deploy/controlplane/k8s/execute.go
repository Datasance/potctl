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

package deployk8scontrolplane

import (
	"fmt"

	"github.com/datasance/potctl/internal/config"
	"github.com/datasance/potctl/internal/execute"
	rsc "github.com/datasance/potctl/internal/resource"
	"github.com/datasance/potctl/pkg/iofog/install"
	"github.com/datasance/potctl/pkg/util"
)

type Options struct {
	Namespace string
	Yaml      []byte
	Name      string
}

type kubernetesControlPlaneExecutor struct {
	controlPlane *rsc.KubernetesControlPlane
	namespace    string
	name         string
}

func (exe kubernetesControlPlaneExecutor) Execute() (err error) {
	util.SpinStart(fmt.Sprintf("Deploying controlplane %s", exe.GetName()))
	if err := exe.executeInstall(); err != nil {
		return err
	}

	// Update config
	ns, err := config.GetNamespace(exe.namespace)
	if err != nil {
		return
	}
	ns.SetControlPlane(exe.controlPlane)
	return config.Flush()
}

func (exe kubernetesControlPlaneExecutor) GetName() string {
	return exe.name
}

func newControlPlaneExecutor(namespace, name string, controlPlane *rsc.KubernetesControlPlane) execute.Executor {
	return kubernetesControlPlaneExecutor{
		namespace:    namespace,
		controlPlane: controlPlane,
		name:         name,
	}
}

func NewExecutor(opt Options) (exe execute.Executor, err error) {
	// Check the namespace exists
	_, err = config.GetNamespace(opt.Namespace)
	if err != nil {
		return
	}

	// Read the input file
	controlPlane, err := rsc.UnmarshallKubernetesControlPlane(opt.Yaml)
	if err != nil {
		return
	}
	if err := validate(&controlPlane); err != nil {
		return nil, err
	}

	return newControlPlaneExecutor(opt.Namespace, opt.Name, &controlPlane), nil
}

func (exe *kubernetesControlPlaneExecutor) executeInstall() (err error) {

	// Get Kubernetes deployer
	installer, err := install.NewKubernetes(exe.controlPlane.KubeConfig, exe.namespace, exe.name)
	if err != nil {
		return
	}

	// Configure deploy
	installer.SetOperatorImage(exe.controlPlane.Images.Operator)
	instller.SetPullSecret(exe.controlPlane.Images.PullSecret)
	installer.SetPortManagerImage(exe.controlPlane.Images.PortManager)
	installer.SetRouterImage(exe.controlPlane.Images.Router)
	installer.SetProxyImage(exe.controlPlane.Images.Proxy)
	installer.SetControllerImage(exe.controlPlane.Images.Controller)
	installer.SetControllerService(exe.controlPlane.Services.Controller.Type, exe.controlPlane.Services.Controller.Address, exe.controlPlane.Services.Controller.Annotations)
	installer.SetRouterService(exe.controlPlane.Services.Router.Type, exe.controlPlane.Services.Router.Address, exe.controlPlane.Services.Router.Annotations)
	installer.SetProxyService(exe.controlPlane.Services.Proxy.Type, exe.controlPlane.Services.Proxy.Address, exe.controlPlane.Services.Proxy.Annotations)
	installer.SetControllerIngress(exe.controlPlane.Ingresses.Controller.Annotations, exe.controlPlane.Ingresses.Controller.IngressClassName, exe.controlPlane.Ingresses.Controller.Host, exe.controlPlane.Ingresses.Controller.SecretName)
	installer.SetRouterIngress(exe.controlPlane.Ingresses.Router.Address, exe.controlPlane.Ingresses.Router.MessagePort, exe.controlPlane.Ingresses.Router.InteriorPort, exe.controlPlane.Ingresses.Router.EdgePort)
	installer.SetHttpProxyIngress(exe.controlPlane.Ingresses.HTTPProxy.Address)
	installer.SetTcpProxyIngress(exe.controlPlane.Ingresses.TCPProxy.Address)

	replicas := int32(1)
	if exe.controlPlane.Replicas.Controller != 0 {
		replicas = exe.controlPlane.Replicas.Controller
	}
	// Create controller on cluster
	// user := install.IofogUser(exe.controlPlane.IofogUser)
	conf := install.ControllerConfig{
		// User:          user,
		Replicas:      replicas,
		Auth:          install.Auth(exe.controlPlane.Auth),
		Database:      install.Database(exe.controlPlane.Database),
		PidBaseDir:    exe.controlPlane.Controller.PidBaseDir,
		EcnViewerPort: exe.controlPlane.Controller.EcnViewerPort,
		EcnViewerURL:  exe.controlPlane.Controller.EcnViewerURL,
		Https:         exe.controlPlane.Controller.Https,
		SecretName:    exe.controlPlane.Controller.SecretName,
	}
	endpoint, err := installer.CreateControlPlane(&conf)
	if err != nil {
		return
	}

	// Create controller pods for config
	pods, err := installer.GetControllerPods()
	if err != nil {
		return
	}
	for idx := range pods {
		k8sPod := rsc.KubernetesController{
			Endpoint: endpoint,
			PodName:  pods[idx].Name,
			Created:  util.NowUTC(),
		}
		if err := exe.controlPlane.AddController(&k8sPod); err != nil {
			return err
		}
	}

	// Assign control plane endpoint
	exe.controlPlane.Endpoint = endpoint

	return err
}

func validate(controlPlane *rsc.KubernetesControlPlane) (err error) {
	// Validate user
	user := controlPlane.GetUser()
	if user.Email == "" || user.Name == "" || user.Password == "" || user.Surname == "" {
		return util.NewInputError("Control Plane Iofog User must contain non-empty values in email, name, surname, and password fields")
	}
	// Validate auth
	auth := controlPlane.Auth
	if auth.URL == "" || auth.Realm == "" || auth.SSL == "" || auth.RealmKey == "" || auth.ControllerClient == "" || auth.ControllerSecret == "" || auth.ViewerClient == "" {
		return util.NewInputError("Control Plane Auth Config must contain non-empty values in all fields")
	}
	// Validate database
	db := controlPlane.Database
	if db.Host == "" || db.DatabaseName == "" || db.Password == "" || db.Port == 0 || db.User == "" {
		msg := `When you are specifying an external database for the Control Plane, you must provide non-empty values in host, databasename, user, password, and port fields.`
		return util.NewInputError(msg)
	}
	// Validate controller service and ingress
	controllerService := controlPlane.Services.Controller
	controllerIngress := controlPlane.Ingresses.Controller
	if controllerService.Type == "ClusterIP" {
		if controllerIngress.Host == "" || controllerIngress.SecretName == "" {
			return util.NewInputError("When Controller service type is ClusterIP, You must provide Ingress configuration for Controller")
		}
	}
	// Validate router service and ingress
	routerService := controlPlane.Services.Router
	routerIngress := controlPlane.Ingresses.Router
	if routerService.Type == "ClusterIP" {
		if routerIngress.Adress == "" || routerIngress.MessagePort == "" || routerIngress.InteriorPort == "" || routerIngress.EdgePort == "" {
			return util.NewInputError("When Router service type is ClusterIP, You must provide Ingress configuration for Default-Router")
		}
	}
	// Validate proxy services and ingress
	proxyService := controlPlane.Services.Proxy
	httpIngress := controlPlane.Ingresses.HTTPProxy
	tcpIngress := controlPlane.Ingresses.TCPProxy
	if proxyService.Type == "ClusterIP" {
		if httpIngress.Address == "" || tcpIngress.Address == "" {
			return util.NewInputError("When Proxy service type is ClusterIP, You must provide Ingress configuration for HTTP and TCP Proxy")
		}
	}
	return
}
