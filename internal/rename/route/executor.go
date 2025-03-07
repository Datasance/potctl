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

package route

import (
	"fmt"

	clientutil "github.com/datasance/potctl/internal/util/client"
	"github.com/datasance/potctl/pkg/util"
)

func Execute(namespace, name, newName string) error {
	// Init remote resources
	clt, err := clientutil.NewControllerClient(namespace)
	if err != nil {
		return err
	}

	appName, routeName, err := clientutil.ParseFQName(name, "Route")
	if err != nil {
		return err
	}

	route, err := clt.GetRoute(appName, routeName)
	if err != nil {
		return err
	}

	if err := util.IsLowerAlphanumeric("Route", newName); err != nil {
		return err
	}
	util.SpinStart(fmt.Sprintf("Renaming route %s", name))
	route.Name = newName
	// Temporary fix
	route.From = ""
	route.To = ""

	if err := clt.PatchRoute(appName, routeName, &route); err != nil {
		return err
	}

	return err
}
