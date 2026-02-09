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
	clientutil "github.com/datasance/potctl/internal/util/client"
	"github.com/datasance/potctl/pkg/util"
)

type serviceAccountExecutor struct {
	namespace string
	name      string
	filename  string
}

func newServiceAccountExecutor(namespace, name, filename string) *serviceAccountExecutor {
	return &serviceAccountExecutor{
		namespace: namespace,
		name:      name,
		filename:  filename,
	}
}

func (exe *serviceAccountExecutor) GetName() string {
	return exe.name
}

func (exe *serviceAccountExecutor) Execute() error {
	clt, err := clientutil.NewControllerClient(exe.namespace)
	if err != nil {
		return err
	}

	sa, err := clt.GetServiceAccount(exe.name)
	if err != nil {
		return err
	}

	spec := rsc.ServiceAccount{
		Name:    sa.Name,
		RoleRef: convertRoleRef(sa.RoleRef),
	}

	header := config.Header{
		APIVersion: config.LatestAPIVersion,
		Kind:       config.ServiceAccountKind,
		Metadata: config.HeaderMetadata{
			Namespace: exe.namespace,
			Name:      exe.name,
		},
		Spec: spec,
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
