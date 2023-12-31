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

package deletecontroller

import (
	"github.com/datasance/potctl/internal/config"
	"github.com/datasance/potctl/internal/execute"
	rsc "github.com/datasance/potctl/internal/resource"
	"github.com/datasance/potctl/pkg/util"
)

func NewExecutor(namespace, name string) (execute.Executor, error) {
	// Get controller from config
	ns, err := config.GetNamespace(namespace)
	if err != nil {
		return nil, err
	}
	baseControlPlane, err := ns.GetControlPlane()
	if err != nil {
		return nil, err
	}
	switch controlPlane := baseControlPlane.(type) {
	case *rsc.KubernetesControlPlane:
		return nil, util.NewInputError("Cannot delete Kubernetes Controller, delete the Control Plane instead.")
	case *rsc.RemoteControlPlane:
		return NewRemoteExecutor(controlPlane, namespace, name), nil
	case *rsc.LocalControlPlane:
		return NewLocalExecutor(controlPlane, namespace, name), nil
	}

	return nil, util.NewInternalError("Could not determine what kind of Control Plane is in Namespace " + namespace)
}
