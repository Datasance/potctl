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

package upgrade

import (
	"github.com/datasance/potctl/internal/execute"
	"github.com/datasance/potctl/pkg/util"
)

type Options struct {
	ResourceType string
	Namespace    string
	Name         string
}

func NewExecutor(opt Options) (execute.Executor, error) {
	switch opt.ResourceType {
	case "agent":
		return newAgentExecutor(opt), nil
	default:
		return nil, util.NewInputError("Unsupported resource: " + opt.ResourceType)
	}
}
