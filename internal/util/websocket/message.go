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
	"fmt"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

// Message represents a WebSocket message with type and payload
type Message struct {
	Type             uint8  `msgpack:"type"`
	Data             []byte `msgpack:"data"`
	MicroserviceUUID string `msgpack:"microserviceUuid"`
	ExecID           string `msgpack:"execId"`
	Timestamp        int64  `msgpack:"timestamp"`
}

// NewMessage creates a new Message with the given type and payload
func NewMessage(msgType uint8, payload []byte, microserviceUUID string, execID string) *Message {
	return &Message{
		Type:             msgType,
		Data:             payload,
		MicroserviceUUID: microserviceUUID,
		ExecID:           execID,
		Timestamp:        time.Now().UnixMilli(),
	}
}

// Encode encodes the message into MessagePack format
func (m *Message) Encode() ([]byte, error) {
	if len(m.Data) > DefaultMaxPayload {
		return nil, fmt.Errorf("payload size exceeds maximum allowed size of %d bytes", DefaultMaxPayload)
	}

	// Encode using MessagePack
	return msgpack.Marshal(m)
}

// Decode decodes a MessagePack message into a Message struct
func Decode(data []byte) (*Message, error) {
	var msg Message
	if err := msgpack.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to decode MessagePack data: %v", err)
	}
	return &msg, nil
}

// IsControlMessage checks if the message is a control message
func (m *Message) IsControlMessage() bool {
	return m.Type == MessageTypeControl
}

// IsCloseMessage checks if the message is a close message
func (m *Message) IsCloseMessage() bool {
	return m.Type == MessageTypeClose
}

// IsActivationMessage checks if the message is an activation message
func (m *Message) IsActivationMessage() bool {
	return m.Type == MessageTypeActivation
}

// IsStdinMessage checks if the message is a stdin message
func (m *Message) IsStdinMessage() bool {
	return m.Type == MessageTypeStdin
}

// IsStdoutMessage checks if the message is a stdout message
func (m *Message) IsStdoutMessage() bool {
	return m.Type == MessageTypeStdout
}

// IsStderrMessage checks if the message is a stderr message
func (m *Message) IsStderrMessage() bool {
	return m.Type == MessageTypeStderr
}
