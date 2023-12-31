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

package describe

import (
	"github.com/datasance/potctl/internal/config"
	rsc "github.com/datasance/potctl/internal/resource"
	"github.com/datasance/potctl/pkg/util"
)

type controllerExecutor struct {
	namespace string
	name      string
	filename  string
}

func newControllerExecutor(namespace, name, filename string) *controllerExecutor {
	c := &controllerExecutor{}
	c.namespace = namespace
	c.name = name
	c.filename = filename
	return c
}

func (exe *controllerExecutor) GetName() string {
	return exe.name
}

func (exe *controllerExecutor) Execute() error {
	ns, err := config.GetNamespace(exe.namespace)
	if err != nil {
		return err
	}
	controlPlane, err := ns.GetControlPlane()
	if err != nil {
		return err
	}
	baseController, err := controlPlane.GetController(exe.name)
	if err != nil {
		return err
	}

	// Generate header
	var header config.Header
	switch controller := baseController.(type) {
	case *rsc.KubernetesController:
		header = exe.generateControllerHeader(config.KubernetesControllerKind, controller)
	case *rsc.RemoteController:
		header = exe.generateControllerHeader(config.RemoteControllerKind, controller)
	case *rsc.LocalController:
		header = exe.generateControllerHeader(config.LocalControllerKind, controller)
	default:
		return util.NewInternalError("Could not convert Control Plane to dynamic type")
	}

	if exe.filename == "" {
		if err := util.Print(header); err != nil {
			return err
		}
	} else {
		if err := util.FPrint(header, exe.filename); err != nil {
			return err
		}
	}
	return nil
}

func (exe *controllerExecutor) generateControllerHeader(kind config.Kind, controller rsc.Controller) config.Header {
	return config.Header{
		APIVersion: config.LatestAPIVersion,
		Kind:       kind,
		Metadata: config.HeaderMetadata{
			Namespace: exe.namespace,
			Name:      controller.GetName(),
		},
		Spec: controller,
	}
}
