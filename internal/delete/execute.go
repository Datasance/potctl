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

package delete

import (
	"fmt"

	"github.com/datasance/potctl/internal/config"
	deleteagent "github.com/datasance/potctl/internal/delete/agent"
	deleteapplication "github.com/datasance/potctl/internal/delete/application"
	deletecatalogitem "github.com/datasance/potctl/internal/delete/catalogitem"
	deletecertificate "github.com/datasance/potctl/internal/delete/certificate"
	deleteconfigmap "github.com/datasance/potctl/internal/delete/configmap"
	deletecontroller "github.com/datasance/potctl/internal/delete/controller"
	deletek8scontrolplane "github.com/datasance/potctl/internal/delete/controlplane/k8s"
	deletelocalcontrolplane "github.com/datasance/potctl/internal/delete/controlplane/local"
	deleteremotecontrolplane "github.com/datasance/potctl/internal/delete/controlplane/remote"
	deletemicroservice "github.com/datasance/potctl/internal/delete/microservice"
	deleteregistry "github.com/datasance/potctl/internal/delete/registry"
	deleterole "github.com/datasance/potctl/internal/delete/role"
	deleterolebinding "github.com/datasance/potctl/internal/delete/rolebinding"
	deleteroute "github.com/datasance/potctl/internal/delete/route"
	deletesecret "github.com/datasance/potctl/internal/delete/secret"
	deleteservice "github.com/datasance/potctl/internal/delete/service"
	deleteserviceaccount "github.com/datasance/potctl/internal/delete/serviceaccount"
	deletevolume "github.com/datasance/potctl/internal/delete/volume"
	deletevolumemount "github.com/datasance/potctl/internal/delete/volumemount"
	"github.com/datasance/potctl/internal/execute"
	"github.com/datasance/potctl/pkg/util"
)

type Options struct {
	Namespace string
	InputFile string
	Soft      bool
}

var kindOrder = []config.Kind{
	config.ServiceKind,
	config.RouteKind,
	config.MicroserviceKind,
	config.ApplicationKind,
	config.CatalogItemKind,
	config.RegistryKind,
	config.VolumeMountKind,
	config.VolumeKind,
	config.LocalAgentKind,
	config.RemoteAgentKind,
	config.ConfigMapKind,
	config.ServiceAccountKind,
	config.RoleBindingKind,
	config.RoleKind,
	config.CertificateKind,
	config.CertificateAuthorityKind,
	config.SecretKind,
	config.LocalControllerKind,
	config.RemoteControllerKind,
	config.LocalControlPlaneKind,
	config.RemoteControlPlaneKind,
	config.KubernetesControlPlaneKind,
}

var kindHandlers = map[config.Kind]func(*execute.KindHandlerOpt) (execute.Executor, error){
	config.ApplicationKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deleteapplication.NewExecutor(opt.Namespace, opt.Name)
	},
	config.MicroserviceKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deletemicroservice.NewExecutor(opt.Namespace, opt.Name)
	},
	config.KubernetesControlPlaneKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deletek8scontrolplane.NewExecutor(opt.Namespace)
	},
	config.RemoteControlPlaneKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deleteremotecontrolplane.NewExecutor(opt.Namespace)
	},
	config.LocalControlPlaneKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deletelocalcontrolplane.NewExecutor(opt.Namespace)
	},
	config.RemoteControllerKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deletecontroller.NewExecutor(opt.Namespace, opt.Name)
	},
	config.LocalControllerKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deletecontroller.NewExecutor(opt.Namespace, opt.Name)
	},
	config.RemoteAgentKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deleteagent.NewExecutor(opt.Namespace, opt.Name, false, false)
	},
	config.LocalAgentKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deleteagent.NewExecutor(opt.Namespace, opt.Name, false, false)
	},
	config.CatalogItemKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deletecatalogitem.NewExecutor(opt.Namespace, opt.Name)
	},
	config.RegistryKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deleteregistry.NewExecutor(opt.Namespace, opt.Name)
	},
	config.VolumeKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deletevolume.NewExecutor(opt.Namespace, opt.Name)
	},
	config.SecretKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deletesecret.NewExecutor(opt.Namespace, opt.Name)
	},
	config.ConfigMapKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deleteconfigmap.NewExecutor(opt.Namespace, opt.Name)
	},
	config.RoleKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deleterole.NewExecutor(opt.Namespace, opt.Name)
	},
	config.RoleBindingKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deleterolebinding.NewExecutor(opt.Namespace, opt.Name)
	},
	config.ServiceAccountKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deleteserviceaccount.NewExecutor(opt.Namespace, opt.Name)
	},
	config.ServiceKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deleteservice.NewExecutor(opt.Namespace, opt.Name)
	},
	config.RouteKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deleteroute.NewExecutor(opt.Namespace, opt.Name), nil
	},
	config.VolumeMountKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deletevolumemount.NewExecutor(opt.Namespace, opt.Name)
	},
	config.CertificateKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deletecertificate.NewExecutor(opt.Namespace, opt.Name)
	},
	config.CertificateAuthorityKind: func(opt *execute.KindHandlerOpt) (exe execute.Executor, err error) {
		return deletecertificate.NewExecutor(opt.Namespace, opt.Name)
	},
}

func Execute(opt *Options) error {
	executorsMap, err := execute.GetExecutorsFromYAML(opt.InputFile, opt.Namespace, kindHandlers)
	if err != nil {
		return err
	}

	// Microservice, Application, Agent, Controller, ControlPlane
	for idx := range kindOrder {
		if errs := execute.RunExecutors(executorsMap[kindOrder[idx]], fmt.Sprintf("delete %s", kindOrder[idx])); len(errs) > 0 {
			for _, err := range errs {
				if _, ok := err.(*util.NotFoundError); !ok {
					return execute.CoalesceErrors(errs)
				}
				util.PrintNotify(fmt.Sprintf("Warning: %s %s.", kindOrder[idx], err.Error()))
			}
		}
	}

	return nil
}
