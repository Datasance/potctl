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

package pruneagent

import (
	"github.com/datasance/potctl/internal/config"
	"github.com/datasance/potctl/internal/execute"
	rsc "github.com/datasance/potctl/internal/resource"
	clientutil "github.com/datasance/potctl/internal/util/client"
	"github.com/datasance/potctl/pkg/util"
)

type executor struct {
	name        string
	namespace   string
	useDetached bool
}

func NewExecutor(namespace, name string, useDetached bool) execute.Executor {
	return executor{name: name, namespace: namespace, useDetached: useDetached}
}

func (exe executor) GetName() string {
	return exe.name
}

func (exe executor) Execute() error {
	util.SpinStart("Pruning Agent")
	ns, err := config.GetNamespace(exe.namespace)
	if err != nil {
		return err
	}
	// Update local cache based on Controller
	if err := clientutil.SyncAgentInfo(exe.namespace); err != nil {
		return err
	}

	var baseAgent rsc.Agent
	if exe.useDetached {
		baseAgent, err = config.GetDetachedAgent(exe.name)
	} else {
		baseAgent, err = ns.GetAgent(exe.name)
	}
	if err != nil {
		return err
	}
	// Prune Agent
	switch agent := baseAgent.(type) {
	case *rsc.LocalAgent:
		if err := exe.localAgentPrune(); err != nil {
			return err
		}
	case *rsc.RemoteAgent:
		if exe.useDetached {
			if err := exe.remoteDetachedAgentPrune(agent); err != nil {
				return err
			}
		} else {
			if err := exe.remoteAgentPrune(baseAgent); err != nil {
				return err
			}
		}
	}
	return nil
}
