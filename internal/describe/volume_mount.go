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

type volumeMountExecutor struct {
	namespace string
	name      string
	filename  string
}

func newVolumeMountExecutor(namespace, name, filename string) *volumeMountExecutor {
	return &volumeMountExecutor{
		namespace: namespace,
		name:      name,
		filename:  filename,
	}
}

func (exe *volumeMountExecutor) GetName() string {
	return exe.name
}

func (exe *volumeMountExecutor) Execute() error {
	// Init remote resources
	clt, err := clientutil.NewControllerClient(exe.namespace)
	if err != nil {
		return err
	}

	// Get Volume Mount from Controller
	volumeMount, err := clt.GetVolumeMount(exe.name)
	if err != nil {
		return err
	}

	header := config.Header{
		APIVersion: config.LatestAPIVersion,
		Kind:       config.VolumeMountKind,
		Metadata: config.HeaderMetadata{
			Namespace: exe.namespace,
			Name:      exe.name,
		},
		Spec: rsc.VolumeMount{
			Name:          volumeMount.Name,
			UUID:          volumeMount.UUID,
			ConfigMapName: volumeMount.ConfigMapName,
			SecretName:    volumeMount.SecretName,
			Version:       volumeMount.Version,
		},
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
