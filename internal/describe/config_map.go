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

type configMapExecutor struct {
	namespace string
	name      string
	filename  string
}

func newConfigMapExecutor(namespace, name, filename string) *configMapExecutor {
	return &configMapExecutor{
		namespace: namespace,
		name:      name,
		filename:  filename,
	}
}

func (exe *configMapExecutor) GetName() string {
	return exe.name
}

func (exe *configMapExecutor) Execute() error {
	// Init remote resources
	clt, err := clientutil.NewControllerClient(exe.namespace)
	if err != nil {
		return err
	}

	// Get secret from Controller
	configMap, err := clt.GetConfigMap(exe.name)
	if err != nil {
		return err
	}

	header := config.Header{
		APIVersion: config.LatestAPIVersion,
		Kind:       config.ConfigMapKind,
		Metadata: config.HeaderMetadata{
			Namespace: exe.namespace,
			Name:      exe.name,
		},
		Spec: rsc.ConfigMap{
			Immutable: configMap.Immutable,
		},
		Data: configMap.Data,
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
