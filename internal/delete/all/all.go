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

package deleteall

import (
	"github.com/datasance/potctl/internal/config"
	deleteagent "github.com/datasance/potctl/internal/delete/agent"
	deletecontrolplane "github.com/datasance/potctl/internal/delete/controlplane"
	deletevolume "github.com/datasance/potctl/internal/delete/volume"
	"github.com/datasance/potctl/internal/execute"
	clientutil "github.com/datasance/potctl/internal/util/client"
	"github.com/datasance/potctl/pkg/util"
)

func Execute(namespace string, useDetached, force bool) error {
	// Make sure to update config despite failure
	defer config.Flush()

	// Get namespace
	ns, err := config.GetNamespace(namespace)
	if err != nil {
		return err
	}

	// Delete Volumes
	if len(ns.Volumes) > 0 {
		util.SpinStart("Deleting Volumes")
		var executors []execute.Executor
		for _, volume := range ns.Volumes {
			exe, err := deletevolume.NewExecutor(namespace, volume.Name)
			if err != nil {
				return err
			}
			executors = append(executors, exe)
		}
		if err := runExecutors(executors); err != nil {
			return err
		}
	}

	if !useDetached {
		// Delete applications
		util.SpinStart("Deleting Flows")
		clt, err := clientutil.NewControllerClient(namespace)
		if err != nil {
			return err
		}

		flows, err := clt.GetAllFlows()
		if err != nil {
			return err
		}

		for _, flow := range flows.Flows {
			if err := clt.DeleteFlow(flow.ID); err != nil {
				return err
			}
		}
	}

	// Delete Agents
	if len(ns.GetAgents()) > 0 {
		util.SpinStart("Deleting Agents")

		var executors []execute.Executor
		for _, agent := range ns.GetAgents() {
			exe, err := deleteagent.NewExecutor(namespace, agent.GetName(), useDetached, force)
			if err != nil {
				return err
			}
			executors = append(executors, exe)
		}
		if err := runExecutors(executors); err != nil {
			return err
		}
	}

	if !useDetached {
		// Delete Controllers
		util.SpinStart("Deleting Control Plane ")
		exe, err := deletecontrolplane.NewExecutor(namespace)
		if err != nil {
			return err
		}
		if err := exe.Execute(); err != nil {
			return err
		}
	}

	return nil
}

func runExecutors(executors []execute.Executor) error {
	if errs, _ := execute.ForParallel(executors); len(errs) > 0 {
		return execute.CoalesceErrors(errs)
	}
	return nil
}
