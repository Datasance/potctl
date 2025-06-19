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

package movemicroservice

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/datasance/iofog-go-sdk/v3/pkg/apps"
	"github.com/datasance/potctl/internal/describe"
	clientutil "github.com/datasance/potctl/internal/util/client"
	"github.com/datasance/potctl/pkg/util"
	"gopkg.in/yaml.v2"
)

func Execute(namespace, name, agent string) error {
	util.SpinStart(fmt.Sprintf("Moving microservice %s", name))

	// Update local cache based on Controller
	if err := clientutil.SyncAgentInfo(namespace); err != nil {
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

	// destAgent, err := clt.GetAgentByName(agent, false)
	destAgent, err := clt.GetAgentByName(agent)
	if err != nil {
		return err
	}

	// Move
	msvc.AgentUUID = destAgent.UUID

	yamlMsvc, _, _, err := describe.MapClientMicroserviceToDeployMicroservice(msvc, clt)
	if err != nil {
		return err
	}

	file := apps.IofogHeader{
		APIVersion: "datasance.com/v3",
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
