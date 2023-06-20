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
	"github.com/datasance/potctl/v1/internal/execute"
	"github.com/datasance/potctl/v1/pkg/util"
)

type Options struct {
	ResourceType string
	Namespace    string
	Name         string
	KubeConfig   string
	KeyFile      string
	User         string
	Port         int
	UseDetached  bool
}

var multipleResources = map[string]bool{
	"agents":      true,
	"controllers": true,
}

func NewExecutor(opt *Options) (execute.Executor, error) {
	switch opt.ResourceType {
	case "current-namespace":
		return newDefaultNamespaceExecutor(opt), nil
	case "default-namespace":
		return newDefaultNamespaceExecutor(opt), nil
	case "controlplane":
		return newControlPlaneExecutor(opt), nil
	case "controller":
		return newControllerExecutor(opt), nil
	case "agent":
		return newAgentExecutor(opt), nil
	default:
		if _, exists := multipleResources[opt.ResourceType]; !exists {
			return nil, util.NewInputError("Unsupported resource: " + opt.ResourceType)
		}
		return newMultipleExecutor(opt), nil
	}
}
