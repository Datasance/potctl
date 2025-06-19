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

type serviceExecutor struct {
	namespace string
	name      string
	filename  string
}

func newServiceExecutor(namespace, name, filename string) *serviceExecutor {
	return &serviceExecutor{
		namespace: namespace,
		name:      name,
		filename:  filename,
	}
}

func (exe *serviceExecutor) GetName() string {
	return exe.name
}

func (exe *serviceExecutor) Execute() error {
	// Init remote resources
	clt, err := clientutil.NewControllerClient(exe.namespace)
	if err != nil {
		return err
	}

	// Get service from Controller
	service, err := clt.GetService(exe.name)
	if err != nil {
		return err
	}

	// Convert tags to pointer
	var tags *[]string
	if len(service.Tags) > 0 {
		tags = &service.Tags
	}

	header := config.Header{
		APIVersion: config.LatestAPIVersion,
		Kind:       config.ServiceKind,
		Metadata: config.HeaderMetadata{
			Namespace: exe.namespace,
			Name:      exe.name,
			Tags:      tags,
		},
		Spec: rsc.ClusterService{
			Type:            service.Type,
			Resource:        service.Resource,
			TargetPort:      service.TargetPort,
			ServicePort:     service.ServicePort,
			K8sType:         service.K8sType,
			BridgePort:      service.BridgePort,
			DefaultBridge:   service.DefaultBridge,
			ServiceEndpoint: service.ServiceEndpoint,
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
