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

package websocket

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/datasance/potctl/pkg/util"
	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client
type Client struct {
	conn             *websocket.Conn
	microserviceUUID string
	execID           string
	done             chan struct{}
	lastMessageType  int // Track the last message type
	err              error
	errMutex         sync.Mutex
	closeOnce        sync.Once
}

// NewClient creates a new WebSocket client
func NewClient(microserviceUUID string) *Client {
	return &Client{
		microserviceUUID: microserviceUUID,
		done:             make(chan struct{}),
	}
}

// Connect establishes a WebSocket connection to the server
func (c *Client) Connect(url string, headers http.Header) error {
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		HandshakeTimeout: 45 * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
	}

	conn, resp, err := dialer.Dial(url, headers)
	if err != nil {
		// Check if this is a close frame error
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				c.errMutex.Lock()
				c.err = util.NewError(fmt.Sprintf("connection closed by server (code: %d, reason: %s)", closeErr.Code, closeErr.Text))
				c.errMutex.Unlock()
				c.Close()
				return c.err
			}
		}
		// Check if we got a non-200 response
		if resp != nil && resp.StatusCode != http.StatusSwitchingProtocols {
			c.errMutex.Lock()
			c.err = util.NewError(fmt.Sprintf("failed to establish WebSocket connection: server returned status %d", resp.StatusCode))
			c.errMutex.Unlock()
			c.Close()
			return c.err
		}
		c.errMutex.Lock()
		c.err = util.NewError(fmt.Sprintf("failed to establish WebSocket connection: %v", err))
		c.errMutex.Unlock()
		c.Close()
		return c.err
	}

	c.conn = conn
	return nil
}

// SendMessage sends a message to the server
func (c *Client) SendMessage(msg *Message) error {
	if c.conn == nil {
		return util.NewError("not connected")
	}

	// Set session-specific fields
	msg.MicroserviceUUID = c.microserviceUUID
	msg.ExecID = c.execID

	data, err := msg.Encode()
	if err != nil {
		return util.NewError(fmt.Sprintf("failed to encode message: %v", err))
	}

	return c.conn.WriteMessage(websocket.BinaryMessage, data)
}

// ReadMessage reads a message from the server
func (c *Client) ReadMessage() (*Message, error) {
	if c.conn == nil {
		return nil, util.NewError("not connected")
	}

	messageType, data, err := c.conn.ReadMessage()
	if err != nil {
		// Check if this is a close frame
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
			// Extract close code and reason if available
			if closeErr, ok := err.(*websocket.CloseError); ok {
				err = util.NewError(fmt.Sprintf("connection closed by server (code: %d, reason: %s)", closeErr.Code, closeErr.Text))
			}
		} else {
			err = util.NewError(fmt.Sprintf("failed to read message: %v", err))
		}
		// Store the error and close connection
		c.errMutex.Lock()
		c.err = err
		c.errMutex.Unlock()
		c.Close()
		return nil, err
	}

	// Store the message type
	c.lastMessageType = messageType

	// All messages are now MessagePack encoded
	msg, err := Decode(data)
	if err != nil {
		// Store the error and close connection
		c.errMutex.Lock()
		c.err = err
		c.errMutex.Unlock()
		c.Close()
		return nil, err
	}

	// Update session-specific fields
	c.execID = msg.ExecID
	return msg, nil
}

// Close closes the WebSocket connection
func (c *Client) Close() error {
	var closeErr error
	c.closeOnce.Do(func() {
		if c.conn != nil {
			// Send close message before closing
			closeMsg := NewMessage(MessageTypeClose, nil, c.microserviceUUID, c.execID)
			_ = c.SendMessage(closeMsg)
			closeErr = c.conn.Close()
			c.conn = nil
		}
		close(c.done)
	})
	return closeErr
}

// IsConnected checks if the client is connected
func (c *Client) IsConnected() bool {
	return c.conn != nil
}

// GetDone returns the done channel
func (c *Client) GetDone() <-chan struct{} {
	return c.done
}

// SetExecID sets the execution ID for the client
func (c *Client) SetExecID(execID string) {
	c.execID = execID
}

// GetExecID returns the current execution ID
func (c *Client) GetExecID() string {
	return c.execID
}

// GetMicroserviceUUID returns the microservice UUID
func (c *Client) GetMicroserviceUUID() string {
	return c.microserviceUUID
}

// GetError returns the last error that occurred
func (c *Client) GetError() error {
	c.errMutex.Lock()
	defer c.errMutex.Unlock()
	return c.err
}
