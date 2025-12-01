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

package deployremotecontroller

import (
	"github.com/datasance/potctl/internal/config"
	"github.com/datasance/potctl/internal/execute"
	rsc "github.com/datasance/potctl/internal/resource"
	"github.com/datasance/potctl/pkg/iofog/install"
	"github.com/datasance/potctl/pkg/util"
)

type remoteExecutor struct {
	namespace    string
	controlPlane *rsc.RemoteControlPlane
	controller   *rsc.RemoteController
}

type Options struct {
	Namespace string
	Yaml      []byte
	Name      string
}

func NewExecutor(opt Options) (exe execute.Executor, err error) {
	controller, err := rsc.UnmarshallRemoteController(opt.Yaml)
	if err != nil {
		return
	}

	if len(opt.Name) > 0 {
		controller.Name = opt.Name
	}

	// Validate
	if err = Validate(&controller); err != nil {
		return
	}

	// Get the Control Plane
	ns, err := config.GetNamespace(opt.Namespace)
	if err != nil {
		return nil, err
	}
	baseControlPlane, err := ns.GetControlPlane()
	if err != nil {
		return
	}
	controlPlane, ok := baseControlPlane.(*rsc.RemoteControlPlane)
	if !ok {
		err = util.NewError("Could not convert Control Plane to Remote Control Plane")
		return
	}

	return NewExecutorWithoutParsing(opt.Namespace, controlPlane, &controller)
}

func newExecutor(namespace string, controlPlane *rsc.RemoteControlPlane, controller *rsc.RemoteController) *remoteExecutor {
	executor := &remoteExecutor{
		namespace:    namespace,
		controlPlane: controlPlane,
		controller:   controller,
	}

	// Set default values
	executor.setDefaultValues()
	return executor
}

func (exe *remoteExecutor) GetName() string {
	return "Deploy Remote Controller"
}

func NewExecutorWithoutParsing(namespace string, controlPlane *rsc.RemoteControlPlane, controller *rsc.RemoteController) (exe execute.Executor, err error) {
	_, err = config.GetNamespace(namespace)
	if err != nil {
		return
	}

	if err := controller.Sanitize(); err != nil {
		return nil, err
	}

	if err := util.IsLowerAlphanumeric("Controller", controller.GetName()); err != nil {
		return nil, err
	}

	// Instantiate executor
	return newExecutor(namespace, controlPlane, controller), nil
}

func (exe *remoteExecutor) Execute() (err error) {
	if err = exe.controller.ValidateSSH(); err != nil {
		return
	}
	if exe.controller.LogLevel == "" {
		exe.controller.LogLevel = "info"
	}
	// Instantiate deployer
	controllerOptions := &install.ControllerOptions{
		Namespace:       exe.namespace,
		User:            exe.controller.SSH.User,
		Host:            exe.controller.Host,
		Port:            exe.controller.SSH.Port,
		PrivKeyFilename: exe.controller.SSH.KeyFile,
		PidBaseDir:      exe.controller.PidBaseDir,
		EcnViewerPort:   exe.controller.EcnViewerPort,
		EcnViewerURL:    exe.controller.EcnViewerURL,
		LogLevel:        exe.controller.LogLevel,
		Version:         exe.controlPlane.Package.Version,
		Image:           exe.controlPlane.Package.Container.Image,
		// Repo:                exe.controlPlane.Package.Repo,
		// Token:               exe.controlPlane.Package.Token,
		SystemMicroservices: exe.controlPlane.SystemMicroservices,
	}

	// Add HTTPS configuration if present
	if exe.controller.Https != nil {
		controllerOptions.Https = &install.Https{
			Enabled: exe.controller.Https.Enabled,
			CACert:  exe.controller.Https.CACert,
			TLSCert: exe.controller.Https.TLSCert,
			TLSKey:  exe.controller.Https.TLSKey,
		}
	}

	// Add SiteCA configuration if present
	if exe.controller.SiteCA != nil {
		controllerOptions.SiteCA = &install.SiteCertificate{
			TLSCert: exe.controller.SiteCA.TLSCert,
			TLSKey:  exe.controller.SiteCA.TLSKey,
		}
	}

	// Add LocalCA configuration if present
	if exe.controller.LocalCA != nil {
		controllerOptions.LocalCA = &install.SiteCertificate{
			TLSCert: exe.controller.LocalCA.TLSCert,
			TLSKey:  exe.controller.LocalCA.TLSKey,
		}
	}

	deployer, err := install.NewController(controllerOptions)
	if err != nil {
		return err
	}

	// Set custom scripts if provided
	if exe.controller.Scripts != nil {
		if err := deployer.CustomizeProcedures(
			exe.controller.Scripts.Directory,
			&exe.controller.Scripts.ControllerProcedures); err != nil {
			return err
		}
	}

	// Set database configuration
	if exe.controlPlane.Database.Host != "" {
		db := exe.controlPlane.Database
		deployer.SetControllerExternalDatabase(db.Host, db.User, db.Password, db.Provider, db.DatabaseName, db.Port, db.SSL, db.CA)
	}

	if exe.controlPlane.Auth.URL != "" {
		auth := exe.controlPlane.Auth
		deployer.SetControllerAuth(auth.URL, auth.Realm, auth.SSL, auth.RealmKey, auth.ControllerClient, auth.ControllerSecret, auth.ViewerClient)
	}

	// Set events configuration if present
	if exe.controlPlane.Events.AuditEnabled != nil {
		auditEnabled := *exe.controlPlane.Events.AuditEnabled
		captureIpAddress := false
		if exe.controlPlane.Events.CaptureIpAddress != nil {
			captureIpAddress = *exe.controlPlane.Events.CaptureIpAddress
		}
		deployer.SetControllerEvents(
			auditEnabled,
			exe.controlPlane.Events.RetentionDays,
			exe.controlPlane.Events.CleanupInterval,
			captureIpAddress,
		)
	}

	// Deploy Controller
	if err = deployer.Install(); err != nil {
		return
	}
	// Update controller
	useHTTPS := false
	if exe.controller.Https != nil && exe.controller.Https.Enabled != nil && *exe.controller.Https.Enabled {
		useHTTPS = true
	}
	exe.controller.Endpoint, err = util.GetControllerEndpoint(exe.controller.Host, useHTTPS)
	if err != nil {
		return err
	}
	return exe.controlPlane.UpdateController(exe.controller)
}

func (exe *remoteExecutor) setDefaultValues() {
	if exe.controlPlane.SystemMicroservices.Router.X86 == "" {
		exe.controlPlane.SystemMicroservices.Router.X86 = util.GetRouterImage()
	}
	if exe.controlPlane.SystemMicroservices.Router.ARM == "" {
		exe.controlPlane.SystemMicroservices.Router.ARM = util.GetRouterARMImage()
	}
}

func Validate(ctrl rsc.Controller) error {
	if err := util.IsLowerAlphanumeric("Controller", ctrl.GetName()); err != nil {
		return err
	}
	return nil
}
