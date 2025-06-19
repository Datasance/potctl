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

package util

import "fmt"

// Set by linker
var (
	versionNumber = "undefined"
	platform      = "undefined"
	commit        = "undefined"
	date          = "undefined"

	repo = "undefined"

	controllerTag     = "undefined"
	agentTag          = "undefined"
	operatorTag       = "undefined"
	routerTag         = "undefined"
	routerAdaptorTag  = "undefined"
	controllerVersion = "undefined"
	agentVersion      = "undefined"
)

const (
	controllerImage    = "controller"
	agentImage         = "agent"
	operatorImage      = "operator"
	routerAdaptorImage = "router-adaptor"
	routerImage        = "router"
	routerARMImage     = "router"
)

type Version struct {
	VersionNumber string `yaml:"version"`
	Platform      string
	Commit        string
	Date          string
}

func GetVersion() Version {
	return Version{
		VersionNumber: versionNumber,
		Platform:      platform,
		Commit:        commit,
		Date:          date,
	}
}

func GetControllerVersion() string { return controllerVersion }
func GetAgentVersion() string      { return agentVersion }

func GetControllerImage() string {
	return fmt.Sprintf("%s/%s:%s", repo, controllerImage, controllerTag)
}
func GetAgentImage() string     { return fmt.Sprintf("%s/%s:%s", repo, agentImage, agentTag) }
func GetOperatorImage() string  { return fmt.Sprintf("%s/%s:%s", repo, operatorImage, operatorTag) }
func GetRouterImage() string    { return fmt.Sprintf("%s/%s:%s", repo, routerImage, routerTag) }
func GetRouterARMImage() string { return fmt.Sprintf("%s/%s:%s", repo, routerARMImage, routerTag) }
func GetRouterAdaptorImage() string {
	return fmt.Sprintf("%s/%s:%s", repo, routerAdaptorImage, routerAdaptorTag)
}
