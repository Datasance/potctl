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

package deleteremotecontrolplane

import (
	"github.com/datasance/potctl/internal/config"
	deletecontroller "github.com/datasance/potctl/internal/delete/controller"
	"github.com/datasance/potctl/internal/execute"
	rsc "github.com/datasance/potctl/internal/resource"
	"github.com/datasance/potctl/pkg/util"
)

type Executor struct {
	namespace string
}

func NewExecutor(namespace string) (execute.Executor, error) {
	exe := &Executor{
		namespace: namespace,
	}
	return exe, nil
}

// GetName returns application name
func (exe *Executor) GetName() string {
	return "Delete Control Plane"
}

// Execute deletes application by deleting its associated flow
func (exe *Executor) Execute() (err error) {
	// Get Control Plane
	ns, err := config.GetNamespace(exe.namespace)
	if err != nil {
		return err
	}
	baseControlPlane, err := ns.GetControlPlane()
	if err != nil {
		return err
	}
	controlPlane, ok := baseControlPlane.(*rsc.RemoteControlPlane)
	if !ok {
		return util.NewError("Could not Convert Controller to Remote Controller")
	}

	controllers := controlPlane.GetControllers()
	executors := make([]execute.Executor, len(controllers))
	for idx := range controllers {
		controller := controllers[idx]
		exe := deletecontroller.NewRemoteExecutor(controlPlane, exe.namespace, controller.GetName())
		executors[idx] = exe
	}

	if err := runExecutors(executors); err != nil {
		return err
	}

	// Delete Control Plane in config
	ns.DeleteControlPlane()
	return config.Flush()
}

func runExecutors(executors []execute.Executor) error {
	if errs, _ := execute.ForParallel(executors); len(errs) > 0 {
		return execute.CoalesceErrors(errs)
	}
	return nil
}
