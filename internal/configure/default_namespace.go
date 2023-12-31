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
	"github.com/datasance/potctl/pkg/util"
)

type defaultNamespaceExecutor struct {
	name string
}

func newDefaultNamespaceExecutor(opt *Options) *defaultNamespaceExecutor {
	return &defaultNamespaceExecutor{
		name: opt.Name,
	}
}

func (exe *defaultNamespaceExecutor) GetName() string {
	return exe.name
}

func (exe *defaultNamespaceExecutor) Execute() error {
	if exe.name == "" {
		return util.NewInputError("Must specify Namespace")
	}
	return config.SetDefaultNamespace(exe.name)
}
