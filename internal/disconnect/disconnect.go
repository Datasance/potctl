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

package disconnect

import (
	"github.com/datasance/potctl/internal/config"
	"github.com/datasance/potctl/pkg/util"
)

type Options struct {
	Namespace string
}

func Execute(opt *Options) error {
	// Check Namespace exists
	if _, err := config.GetNamespace(opt.Namespace); err != nil {
		if util.IsNotFoundError(err) {
			// Not found, disconnection is idempotent
			return nil
		}
		// Error was not 'not found'
		return err
	}

	if err := config.DeleteNamespace(opt.Namespace); err != nil {
		return err
	}
	return config.Flush()
}
