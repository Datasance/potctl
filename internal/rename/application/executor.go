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

package application

import (
	"fmt"

	"github.com/datasance/potctl/internal/config"
	clientutil "github.com/datasance/potctl/internal/util/client"
	"github.com/datasance/potctl/pkg/util"
	"github.com/datasance/iofog-go-sdk/pkg/client"
)

func Execute(namespace, name, newName string) error {
	if err := util.IsLowerAlphanumeric("Application", newName); err != nil {
		return err
	}
	util.SpinStart(fmt.Sprintf("Renaming Application %s", name))

	// Init remote resources
	clt, err := clientutil.NewControllerClient(namespace)
	if err != nil {
		return err
	}

	flow, err := clt.GetFlowByName(name)
	if err != nil {
		return err
	}

	flow.Name = newName
	_, err = clt.UpdateFlow(&client.FlowUpdateRequest{
		ID:   flow.ID,
		Name: &newName,
	})
	if err != nil {
		return err
	}
	config.Flush()
	return nil
}
