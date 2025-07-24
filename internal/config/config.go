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

package config

import (
	"errors"
	"fmt"
	"os"
	"path"

	rsc "github.com/datasance/potctl/internal/resource"
	"github.com/datasance/potctl/pkg/util"
	homedir "github.com/mitchellh/go-homedir"
	yaml "gopkg.in/yaml.v2"
)

var (
	conf               configuration
	configFolder       string // config directory
	configFilename     string // config file name
	namespaceDirectory string // Path of namespace directory
	namespaces         map[string]*rsc.Namespace
)

const (
	apiVersionGroup      = "datasance.com"
	latestVersion        = "v3"
	LatestAPIVersion     = apiVersionGroup + "/" + latestVersion
	defaultDirname       = ".iofog/" + latestVersion
	namespaceDirname     = "namespaces/"
	defaultFilename      = "config.yaml"
	configV3             = "potctl/v3"
	CurrentConfigVersion = configV3
	detachedNamespace    = "_detached"
)

// Init initializes config, namespace and unmarshalls the files
func Init(configFolderArg string) {
	namespaces = make(map[string]*rsc.Namespace)

	var err error
	configFolder, err = util.FormatPath(configFolderArg)
	util.Check(err)

	if configFolder == "" {
		// Find home directory.
		home, err := homedir.Dir()
		util.Check(err)
		configFolder = path.Join(home, defaultDirname)
	} else {
		dirInfo, err := os.Stat(configFolder)
		util.Check(err)
		if !dirInfo.IsDir() {
			util.Check(util.NewInputError(fmt.Sprintf("The config folder %s is not a valid directory", configFolder)))
		}
	}

	// Set default filename if necessary
	filename := path.Join(configFolder, defaultFilename)
	configFilename = filename
	namespaceDirectory = path.Join(configFolder, namespaceDirname)

	// Check config file already exists
	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		err = os.MkdirAll(configFolder, 0755)
		util.Check(err)

		// Create default config file
		conf.DefaultNamespace = "default"
		err = flushShared()
		util.Check(err)
	}

	// Unmarshall the config file
	confHeader := potctlConfig{}
	err = util.UnmarshalYAML(configFilename, &confHeader)
	util.Check(err)

	conf, err = getConfigFromHeader(&confHeader)
	util.Check(err)

	// Check namespace dir exists
	initNamespaces := []string{"default", detachedNamespace}
	flush := false
	for _, initNamespace := range initNamespaces {
		nsFile := getNamespaceFile(initNamespace)
		if _, err := os.Stat(nsFile); os.IsNotExist(err) {
			flush = true
			err = os.MkdirAll(namespaceDirectory, 0755)
			util.Check(err)

			// Create default namespace file
			if err = AddNamespace(initNamespace, util.NowUTC()); err != nil {
				util.Check(errors.New("Could not initialize " + initNamespace + " configuration"))
			}
		}
	}
	if flush {
		err = flushNamespaces()
		util.Check(err)
	}
}

// getNamespaceFile helper function that returns the full path to a namespace file
func getNamespaceFile(name string) string {
	return path.Join(namespaceDirectory, name+".yaml")
}

func getConfigFromHeader(header *potctlConfig) (conf configuration, err error) {
	switch header.APIVersion {
	case CurrentConfigVersion:
		// All good
		break
		// Example for further maintenance
		// case PreviousConfigVersion
		// 	updateFromPreviousVersion()
		// 	break

	}
	bytes, err := yaml.Marshal(header.Spec)
	if err != nil {
		return
	}
	if err = yaml.Unmarshal(bytes, &conf); err != nil {
		return
	}
	return conf, err
}

func getNamespaceFromHeader(header *potctlNamespace) (ns *rsc.Namespace, err error) {
	// Check header not supported
	switch header.APIVersion {
	case CurrentConfigVersion:
		// All good
		break

	}
	// Unmarshal Namespace spec
	bytes, err := yaml.Marshal(header.Spec)
	if err != nil {
		return
	}
	ns = new(rsc.Namespace)
	if err = yaml.Unmarshal(bytes, &ns); err != nil {
		return
	}
	return
}

func getConfigYAMLFile(conf configuration) ([]byte, error) {
	confHeader := potctlConfig{
		Header: Header{
			Kind:       potctlConfigKind,
			APIVersion: CurrentConfigVersion,
			Spec:       conf,
		},
	}

	return yaml.Marshal(confHeader)
}

func getNamespaceYAMLFile(ns *rsc.Namespace) ([]byte, error) {
	namespaceHeader := potctlNamespace{
		Header{
			Kind:       potctlNamespaceKind,
			APIVersion: CurrentConfigVersion,
			Metadata: HeaderMetadata{
				Name: ns.Name,
			},
			Spec: ns,
		},
	}
	return yaml.Marshal(namespaceHeader)
}

func flushNamespaces() error {
	for _, ns := range namespaces {
		// Marshal the runtime data
		marshal, err := getNamespaceYAMLFile(ns)
		if err != nil {
			return err
		}
		// Overwrite the file
		err = os.WriteFile(getNamespaceFile(ns.Name), marshal, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func flushShared() error {
	// Marshal the runtime data
	marshal, err := getConfigYAMLFile(conf)
	if err != nil {
		return nil
	}
	// Overwrite the file
	err = os.WriteFile(configFilename, marshal, 0644)
	if err != nil {
		return nil
	}
	return nil
}

// Flush will write namespace files to disk
func Flush() error {
	return flushNamespaces()
}

func ValidateHeader(header *Header) error {
	if header.APIVersion != LatestAPIVersion {
		return util.NewInputError(fmt.Sprintf("Unsupported YAML API version %s.\nPlease use version %s. See https://docs.datasance.com for specification details.", header.APIVersion, LatestAPIVersion))
	}
	return nil
}
