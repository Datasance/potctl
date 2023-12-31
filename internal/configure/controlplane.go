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

package configure

import (
	"github.com/datasance/potctl/internal/config"
	rsc "github.com/datasance/potctl/internal/resource"
	"github.com/datasance/potctl/pkg/util"
)

type kubernetesConfig struct {
	kubeConfig string
}

type controlPlaneExecutor struct {
	namespace        string
	kubernetesConfig kubernetesConfig
	name             string
	remoteConfig     remoteConfig
}

func newControlPlaneExecutor(opt *Options) *controlPlaneExecutor {
	return &controlPlaneExecutor{
		namespace: opt.Namespace,
		name:      opt.Name,
		remoteConfig: remoteConfig{
			keyFile: opt.KeyFile,
			user:    opt.User,
			port:    opt.Port,
		},
		kubernetesConfig: kubernetesConfig{
			kubeConfig: opt.KubeConfig,
		},
	}
}

func (exe *controlPlaneExecutor) GetName() string {
	return exe.name
}

func (exe *controlPlaneExecutor) Execute() error {
	// Get config
	ns, err := config.GetNamespace(exe.namespace)
	if err != nil {
		return err
	}
	baseControlPlane, err := ns.GetControlPlane()
	if err != nil {
		return err
	}

	switch controlPlane := baseControlPlane.(type) {
	case *rsc.RemoteControlPlane:
		return util.NewInputError("Cannot configure Remote Control Plane as if it is a Kubernetes Control Plane")

	case *rsc.KubernetesControlPlane:
		if err := exe.kubernetesConfigure(controlPlane); err != nil {
			return err
		}

	case *rsc.LocalControlPlane:
		return util.NewInputError("Cannot configure a Local Control Plane")
	}

	ns.SetControlPlane(baseControlPlane)
	return config.Flush()
}

func (exe *controlPlaneExecutor) kubernetesConfigure(controlPlane *rsc.KubernetesControlPlane) (err error) {
	// Error if remoteConfig is passed
	if (remoteConfig{}) != exe.remoteConfig {
		return util.NewInputError("Cannot edit remote config of a Kubernetes Control Plane")
	}

	if exe.kubernetesConfig.kubeConfig != "" {
		controlPlane.KubeConfig = exe.kubernetesConfig.kubeConfig
	}

	return controlPlane.Sanitize()
}
