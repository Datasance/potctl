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

package exec

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/datasance/iofog-go-sdk/v3/pkg/client"
	clientutil "github.com/datasance/potctl/internal/util/client"
	"github.com/datasance/potctl/internal/util/terminal"
	"github.com/datasance/potctl/internal/util/websocket"
	"github.com/datasance/potctl/pkg/util"
)

type microserviceExecutor struct {
	namespace string
	name      string
	client    *client.Client
	msvc      *client.MicroserviceInfo
}

func newMicroserviceExecutor(namespace, name string) *microserviceExecutor {
	a := &microserviceExecutor{}
	a.namespace = namespace
	a.name = name
	return a
}

func (exe *microserviceExecutor) GetName() string {
	return exe.name
}

func (exe *microserviceExecutor) Execute() error {
	util.SpinStart("Connecting Exec Session to Microservice")

	// Init client
	clt, err := clientutil.NewControllerClient(exe.namespace)
	if err != nil {
		return err
	}

	appName, msvcName, err := clientutil.ParseFQName(exe.name, "Microservice")
	if err != nil {
		return err
	}

	exe.msvc, err = clt.GetMicroserviceByName(appName, msvcName)
	if err != nil {
		// Check if error indicates application not found
		if strings.Contains(err.Error(), "Invalid application id") {
			// Try system application
			exe.msvc, err = clt.GetSystemMicroserviceByName(appName, msvcName)
			if err != nil {
				return err
			}
		} else {
			// Return other types of errors
			return err
		}
	}

	// Create WebSocket client
	wsClient := websocket.NewClient(exe.msvc.UUID)

	// Get controller endpoint
	controllerURL := clt.GetBaseURL()
	// Convert http(s):// to ws(s)://
	wsURL := strings.Replace(controllerURL, "http://", "ws://", 1)
	wsURL = strings.Replace(wsURL, "https://", "wss://", 1)
	wsURL = fmt.Sprintf("%s/microservices/exec/%s", wsURL, exe.msvc.UUID)

	// Set up headers
	headers := http.Header{}
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", clt.GetAccessToken()))
	util.SpinHandlePrompt()
	// Connect to WebSocket
	if err := wsClient.Connect(wsURL, headers); err != nil {
		util.SpinHandlePromptComplete()
		return util.NewError(fmt.Sprintf("failed to connect to WebSocket: %v", err))
	}

	// Create and start terminal
	term := terminal.NewTerminal(wsClient)

	// Check for initial connection error
	if err := wsClient.GetError(); err != nil {
		util.SpinHandlePromptComplete()
		formattedErr := formatWebSocketError(err)
		return util.NewError(formattedErr)
	}

	if err := term.Start(); err != nil {
		util.SpinHandlePromptComplete()
		formattedErr := formatWebSocketError(err)
		return util.NewError(formattedErr)
	}

	// Wait for terminal to finish
	<-wsClient.GetDone()
	msg := fmt.Sprintf("Successfully closed Microservice %s Exec Session", exe.name)
	util.PrintSuccess(msg)

	// Check if there was an error
	if err := wsClient.GetError(); err != nil {
		formattedErr := formatWebSocketError(err)
		return util.NewError(formattedErr)
	}

	return nil
}
