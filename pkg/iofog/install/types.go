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

package install

type IofogUser struct {
	Name            string
	Surname         string
	Email           string
	Password        string
	SubscriptionKey string
}

type Auth struct {
	URL              string
	Realm            string
	SSL              string
	RealmKey         string
	ControllerClient string
	ControllerSecret string
	ViewerClient     string
}

type Database struct {
	Provider     string
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
}

type Pod struct {
	Name   string
	Status string
}

type ControllerConfig struct {
	// User          IofogUser
	Replicas      int32
	Database      Database
	PidBaseDir    string
	EcnViewerPort int
	EcnViewerURL  string
	Auth          Auth
	Https         *bool
	SecretName    string
}
