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

package deleteapplication

import (
	"github.com/datasance/potctl/v1/internal/execute"
	clientutil "github.com/datasance/potctl/v1/internal/util/client"
	"github.com/datasance/potctl/v1/pkg/util"
	"github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client"
)

type Executor struct {
	namespace string
	name      string
	client    *client.Client
	flow      *client.FlowInfo
}

func NewExecutor(namespace, name string) (execute.Executor, error) {
	exe := &Executor{
		namespace: namespace,
		name:      name,
	}

	return exe, nil
}

// GetName returns application name
func (exe *Executor) GetName() string {
	return exe.name
}

func (exe *Executor) init() (err error) {
	exe.client, err = clientutil.NewControllerClient(exe.namespace)
	if err != nil {
		return
	}
	return
}

// Execute deletes application by deleting its associated flow
func (exe *Executor) Execute() (err error) {
	util.SpinStart("Deleting Application")
	if err := exe.init(); err != nil {
		return err
	}

	err = exe.client.DeleteApplication(exe.name)
	// If notfound error, try legacy
	if _, ok := err.(*client.NotFoundError); ok {
		return exe.deleteLegacy()
	}
	return err
}
