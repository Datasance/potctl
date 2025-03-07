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

package iofog

import "github.com/datasance/iofog-go-sdk/v3/pkg/client"

// String and numeric values of TCP ports used accross ioFog
const (
	ControllerPort       = client.ControllerPort
	ControllerPortString = client.ControllerPortString

	ControllerHostECNViewerPort       = 8008
	ControllerHostECNViewerPortString = "8008"

	DefaultHTTPPort       = 8008
	DefaultHTTPPortString = "8008"

	VanillaRemoteAgentName string = "0-controlplane"
	VanillaRouterAgentName string = client.DefaultRouterName
)
