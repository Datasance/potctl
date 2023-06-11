/*
 *  *******************************************************************************
 *  * Copyright (c) 2020 Edgeworx, Inc.
 *  *
 *  * This program and the accompanying materials are made available under the
 *  * terms of the Eclipse Public License v. 2.0 which is available at
 *  * http://www.eclipse.org/legal/epl-2.0
 *  *
 *  * SPDX-License-Identifier: EPL-2.0
 *  *******************************************************************************
 *
 */

package microservice

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/datasance/potctl/internal/describe"
	clientutil "github.com/datasance/potctl/internal/util/client"
	"github.com/datasance/potctl/pkg/util"
	"github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/apps"
	"gopkg.in/yaml.v2"
)

func Execute(namespace, name, newName string) error {
	if err := util.IsLowerAlphanumeric("Microservice", newName); err != nil {
		return err
	}
	// Init remote resources
	clt, err := clientutil.NewControllerClient(namespace)
	if err != nil {
		return err
	}

	appName, msvcName, err := clientutil.ParseFQName(name, "Microservice")
	if err != nil {
		return err
	}

	msvc, err := clt.GetMicroserviceByName(appName, msvcName)
	if err != nil {
		return err
	}

	util.SpinStart(fmt.Sprintf("Renaming microservice %s", name))

	// Move
	msvc.Name = newName

	yamlMsvc, err := describe.MapClientMicroserviceToDeployMicroservice(msvc, clt)
	if err != nil {
		return err
	}

	file := apps.IofogHeader{
		APIVersion: "iofog.org/v3",
		Kind:       apps.MicroserviceKind,
		Metadata: apps.HeaderMetadata{
			Name: strings.Join([]string{msvc.Application, msvc.Name}, "/"),
		},
		Spec: yamlMsvc,
	}
	yamlBytes, err := yaml.Marshal(file)
	if err != nil {
		return err
	}

	_, err = clt.UpdateMicroserviceFromYAML(msvc.UUID, bytes.NewReader(yamlBytes))

	return err
}
