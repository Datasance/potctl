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

package get

import (
	"time"
	"fmt"
	"github.com/datasance/potctl/internal/config"
	rsc "github.com/datasance/potctl/internal/resource"
	clientutil "github.com/datasance/potctl/internal/util/client"
	"github.com/datasance/potctl/pkg/iofog/install"
	"github.com/datasance/potctl/pkg/util"
	"github.com/datasance/iofog-go-sdk/v3/pkg/client"
)

type controllerExecutor struct {
	namespace string
}

func newControllerExecutor(namespace string) *controllerExecutor {
	c := &controllerExecutor{}
	c.namespace = namespace
	return c
}

func (exe *controllerExecutor) GetName() string {
	return ""
}

func (exe *controllerExecutor) Execute() error {
	table, err := generateControllerOutput(exe.namespace)
	if err != nil {
		return err
	}
	printNamespace(exe.namespace)
	return print(table)
}

func generateControllerOutput(namespace string) (table [][]string, err error) {
	// Get controller config details
	ns, err := config.GetNamespace(namespace)
	if err != nil {
		return
	}

	podStatuses := []string{}
	// Handle k8s
	baseControlPlane, err := ns.GetControlPlane()
	if err != nil {
		if rsc.IsNoControlPlaneError(err) {
			err = nil
		} else {
			return
		}
	}
	if controlPlane, ok := baseControlPlane.(*rsc.KubernetesControlPlane); ok {
		if err = updateControllerPods(controlPlane, namespace); err != nil {
			return
		}
		ns.SetControlPlane(controlPlane)
		if err = config.Flush(); err != nil {
			return
		}
		for idx := range controlPlane.ControllerPods {
			podStatuses = append(podStatuses, controlPlane.ControllerPods[idx].Status)
		}
	}

	// Handle remote and local
	controllers := ns.GetControllers()


	controlPlaneForUser, err := ns.GetControlPlane()
	if err != nil {
		return nil, err
	}
	
	user := controlPlaneForUser.GetUser()

	// Generate table and headers
	table = make([][]string, len(controllers)+1)
	headers := []string{"CONTROLLER", "STATUS", "AGE", "UPTIME", "VERSION", "ADDR", "PORT", "SUBSCRIPTION EXPIRY DATE", "MAX AGENT SEATS"}
	table[0] = append(table[0], headers...)

	// Populate rows
	for idx, ctrlConfig := range controllers {
		// Instantiate connection to controller
		ctrl, err := clientutil.NewControllerClient(namespace)
		if err != nil {
			return table, err
		}

		// Ping status
		ctrlStatus, err := ctrl.GetStatus()
		uptime := "-"
		status := "Failing"
		if err == nil {
			uptime = util.FormatDuration(time.Duration(int64(ctrlStatus.UptimeSeconds)) * time.Second)
			status = ctrlStatus.Status
		}
		// Handle k8s pod statuses
		if len(podStatuses) != 0 && idx < len(podStatuses) {
			status = podStatuses[idx]
		}

		// Get age
		age := "-"
		if ctrlConfig.GetCreatedTime() != "" {
			age, _ = util.ElapsedUTC(ctrlConfig.GetCreatedTime(), util.NowUTC())
		}
		
		addr, port := getAddressAndPort(ctrlConfig.GetEndpoint(), client.ControllerPortString)

		endpoint, err := controlPlaneForUser.GetEndpoint()
		if err != nil {
			return nil, err
		}

		baseURL, err := util.GetBaseURL(endpoint)
		if err != nil {
			return nil, err
		}

		ctrl, err, subscriptionKey := client.RefreshUserSubscriptionKey(client.Options{BaseURL: baseURL}, user.Email, user.GetRawPassword())
		if err != nil {
			fmt.Println("")
		}

		if ctrl == nil {
			fmt.Println("")
		}

		if subscriptionKey != "" {
			fmt.Println("Subscription Key will be updated from controller")
			user.SubscriptionKey = subscriptionKey
		}

		expiryDate, agentSeats, err := util.GetEntitlementDatasance(user.SubscriptionKey, namespace, user.Email)


		if err != nil {
			return nil, err
		}

		row := []string{
			ctrlConfig.GetName(),
			status,
			age,
			uptime,
			ctrlStatus.Versions.Controller,
			addr,
			port,
			expiryDate,
			agentSeats,
		}
		table[idx+1] = append(table[idx+1], row...)
	}

	return table, err
}

func updateControllerPods(controlPlane *rsc.KubernetesControlPlane, namespace string) (err error) {
	// Clear existing
	controlPlane.ControllerPods = []rsc.KubernetesController{}
	// Get pods
	installer, err := install.NewKubernetes(controlPlane.KubeConfig, namespace)
	if err != nil {
		return
	}
	pods, err := installer.GetControllerPods()
	if err != nil {
		return
	}
	// Add pods
	for idx := range pods {
		k8sPod := rsc.KubernetesController{
			Endpoint: controlPlane.Endpoint,
			PodName:  pods[idx].Name,
			Status:   pods[idx].Status,
		}
		if err := controlPlane.AddController(&k8sPod); err != nil {
			return err
		}
	}
	return
}
