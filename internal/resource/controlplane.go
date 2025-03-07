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

package resource

type ControlPlane interface {
	GetUser() IofogUser
	// UpdateUserSubscriptionKey(string) IofogUser
	UpdateUserTokens(string, string) IofogUser
	GetControllers() []Controller
	GetController(string) (Controller, error)
	GetEndpoint() (string, error)
	UpdateController(Controller) error
	AddController(Controller) error
	DeleteController(string) error
	Sanitize() error
	Clone() ControlPlane
}
