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
	// "fmt"

	jsoniter "github.com/json-iterator/go"

	apps "github.com/datasance/iofog-go-sdk/v3/pkg/apps"
	"github.com/datasance/iofog-go-sdk/v3/pkg/client"
	// "github.com/datasance/potctl/pkg/iofog"
	// "github.com/datasance/potctl/pkg/util"
)

func MapClientMicroserviceToDeployMicroservice(msvc *client.MicroserviceInfo, clt *client.Client) (*apps.Microservice, *apps.MicroserviceStatusInfo, *apps.MicroserviceExecStatusInfo, error) {
	agent, err := clt.GetAgentByID(msvc.AgentUUID)
	if err != nil {
		return nil, nil, nil, err
	}
	var catalogItem *client.CatalogItemInfo
	if msvc.CatalogItemID != 0 {
		catalogItem, err = clt.GetCatalogItem(msvc.CatalogItemID)
		if err != nil {
			if httpErr, ok := err.(*client.HTTPError); ok && httpErr.Code == 404 {
				catalogItem = nil
			} else {
				return nil, nil, nil, err
			}
		}
	}

	applicationName := msvc.Application
	if msvc.Application == "" {
		if msvc.FlowID > 0 {
			// Legacy
			flow, err := clt.GetFlowByID(msvc.FlowID)
			if err != nil {
				return nil, nil, nil, err
			}
			applicationName = flow.Name
		}
	}

	return constructMicroservice(msvc, agent.Name, applicationName, catalogItem)
}

// func MapClientMicroserviceStatusToDeployMicroserviceStatus(msvc *client.MicroserviceInfo, clt *client.Client) (*apps.MicroserviceStatusInfo, *apps.MicroserviceExecStatusInfo, error) {
// 	msvcStatus := new(apps.MicroserviceStatusInfo)
// 	msvcStatus.Status = msvc.Status.Status
// 	msvcStatus.StartTime = msvc.Status.StartTime
// 	msvcStatus.OperatingDuration = msvc.Status.OperatingDuration
// 	msvcStatus.MemoryUsage = msvc.Status.MemoryUsage
// 	msvcStatus.CPUUsage = msvc.Status.CPUUsage
// 	msvcStatus.ContainerID = msvc.Status.ContainerID
// 	msvcStatus.Percentage = msvc.Status.Percentage
// 	msvcStatus.IPAddress = msvc.Status.IPAddress
// 	msvcStatus.ErrorMessage = msvc.Status.ErrorMessage
// 	msvcStatus.ExecSessionIDs = msvc.Status.ExecSessionIDs

// 	msvcExecStatus := new(apps.MicroserviceExecStatusInfo)
// 	msvcExecStatus.Status = msvc.ExecStatus.Status
// 	msvcExecStatus.ExecSessionID = msvc.ExecStatus.ExecSessionID

// 	return msvcStatus, msvcExecStatus, nil
// }

func constructMicroservice(msvcInfo *client.MicroserviceInfo, agentName, appName string, catalogItem *client.CatalogItemInfo) (msvc *apps.Microservice, status *apps.MicroserviceStatusInfo, execStatus *apps.MicroserviceExecStatusInfo, err error) {
	msvc = new(apps.Microservice)
	msvc.UUID = msvcInfo.UUID
	msvc.Name = msvcInfo.Name
	msvc.Agent = apps.MicroserviceAgent{
		Name: agentName,
	}
	var armImage, x86Image string
	var msvcImages []client.CatalogImage
	if catalogItem != nil {
		msvcImages = catalogItem.Images
	} else {
		msvcImages = msvcInfo.Images
	}
	for _, image := range msvcImages {
		switch client.AgentTypeIDAgentTypeDict[image.AgentTypeID] {
		case "x86":
			x86Image = image.ContainerImage
		case "arm":
			armImage = image.ContainerImage
		default:
		}
	}
	var registryID int
	var imgArray []client.CatalogImage
	if catalogItem != nil {
		registryID = catalogItem.RegistryID
		imgArray = catalogItem.Images
	} else {
		registryID = msvcInfo.RegistryID
		imgArray = msvcInfo.Images
	}
	images := apps.MicroserviceImages{
		CatalogID: msvcInfo.CatalogItemID,
		X86:       x86Image,
		ARM:       armImage,
		Registry:  client.RegistryTypeIDRegistryTypeDict[registryID],
	}
	for _, img := range imgArray {
		if img.AgentTypeID == 1 {
			images.X86 = img.ContainerImage
		} else if img.AgentTypeID == 2 {
			images.ARM = img.ContainerImage
		}
	}
	volumes := mapVolumes(msvcInfo.Volumes)
	envs := mapEnvs(msvcInfo.Env)
	extraHosts := mapExtraHosts(msvcInfo.ExtraHosts)
	msvc.Images = &images
	jsonConfig := make(map[string]interface{})
	if err := jsoniter.Unmarshal([]byte(msvcInfo.Config), &jsonConfig); err != nil {
		return msvc, nil, nil, err
	}
	jsonAnnotations := make(map[string]interface{})
	if err := jsoniter.Unmarshal([]byte(msvcInfo.Annotations), &jsonAnnotations); err != nil {
		return msvc, nil, nil, err
	}
	msvc.Config = jsonConfig
	msvc.Container.Annotations = jsonAnnotations
	msvc.Container.RootHostAccess = msvcInfo.RootHostAccess
	msvc.Container.PidMode = msvcInfo.PidMode
	msvc.Container.IpcMode = msvcInfo.IpcMode
	msvc.Container.Runtime = msvcInfo.Runtime
	msvc.Container.Platform = msvcInfo.Platform
	msvc.Container.RunAsUser = msvcInfo.RunAsUser
	msvc.Container.CdiDevices = msvcInfo.CdiDevices
	msvc.Container.CapAdd = msvcInfo.CapAdd
	msvc.Container.CapDrop = msvcInfo.CapDrop
	msvc.Container.Commands = msvcInfo.Commands
	msvc.Container.Ports = mapPorts(msvcInfo.Ports)
	msvc.Container.Volumes = &volumes
	msvc.Container.Env = &envs
	msvc.Container.ExtraHosts = &extraHosts
	msvc.MsRoutes = apps.MsRoutes{
		PubTags: msvcInfo.PubTags,
		SubTags: msvcInfo.SubTags,
	}
	msvc.Schedule = msvcInfo.Schedule
	msvc.Application = &appName
	status = new(apps.MicroserviceStatusInfo)

	status.Status = msvcInfo.Status.Status
	status.StartTime = msvcInfo.Status.StartTime
	status.OperatingDuration = msvcInfo.Status.OperatingDuration
	status.MemoryUsage = msvcInfo.Status.MemoryUsage
	status.CPUUsage = msvcInfo.Status.CPUUsage
	status.ContainerID = msvcInfo.Status.ContainerID
	status.Percentage = msvcInfo.Status.Percentage
	status.ErrorMessage = msvcInfo.Status.ErrorMessage
	status.IPAddress = msvcInfo.Status.IPAddress
	status.ExecSessionIDs = msvcInfo.Status.ExecSessionIDs
	execStatus = new(apps.MicroserviceExecStatusInfo)
	execStatus.Status = msvcInfo.ExecStatus.Status
	execStatus.ExecSessionID = msvcInfo.ExecStatus.ExecSessionID

	return msvc, status, execStatus, err
}

func mapPort(in *client.MicroservicePortMappingInfo) (out *apps.MicroservicePortMapping) {
	if in == nil {
		return nil
	}
	return &apps.MicroservicePortMapping{
		Internal: in.Internal,
		External: in.External,
		Protocol: in.Protocol,
	}
}

func mapPorts(in []client.MicroservicePortMappingInfo) (out []apps.MicroservicePortMapping) {
	for idx := range in {
		port := mapPort(&in[idx])
		if port != nil {
			out = append(out, *port)
		}
	}
	return
}

func mapVolumes(in []client.MicroserviceVolumeMappingInfo) (out []apps.MicroserviceVolumeMapping) {
	for _, vol := range in {
		out = append(out, apps.MicroserviceVolumeMapping(vol))
	}
	return
}

func mapEnvs(in []client.MicroserviceEnvironmentInfo) (out []apps.MicroserviceEnvironment) {
	for _, env := range in {
		out = append(out, apps.MicroserviceEnvironment(env))
	}
	return
}

func mapExtraHosts(in []client.MicroserviceExtraHost) (out []apps.MicroserviceExtraHost) {
	for _, eH := range in {
		out = append(out, apps.MicroserviceExtraHost(eH))
	}
	return
}
