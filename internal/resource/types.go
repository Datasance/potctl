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

import (
	"time"

	"github.com/datasance/iofog-go-sdk/v3/pkg/apps"
	"github.com/datasance/iofog-go-sdk/v3/pkg/client"
)

type Route = apps.Route
type Microservice = apps.Microservice
type Application = apps.Application
type ApplicationTemplate = apps.ApplicationTemplate

type Container struct {
	Image       string      `yaml:"image,omitempty"`
	Credentials Credentials `yaml:"credentials,omitempty"` // Optional credentials if needed to pull images
}

type RemoteContainer struct {
	Image string `yaml:"image,omitempty"`
	// Repo        string      `yaml:"repo,omitempty"`
	// Credentials Credentials `yaml:"credentials,omitempty"` // Optional credentials if needed to pull images
}

type Package struct {
	Version   string `yaml:"version,omitempty"`
	Container RemoteContainer
	// Repo    string `yaml:"repo,omitempty"`
	// Token   string `yaml:"token,omitempty"`
}

type SSH struct {
	User    string `yaml:"user,omitempty"`
	Port    int    `yaml:"port,omitempty"`
	KeyFile string `yaml:"keyFile,omitempty"`
}

type KubeImages struct {
	PullSecret    string `yaml:"pullSecret,omitempty"`
	Controller    string `yaml:"controller,omitempty"`
	Operator      string `yaml:"operator,omitempty"`
	Router        string `yaml:"router,omitempty"`
	RouterAdaptor string `yaml:"routerAdaptor,omitempty"`
}

type Services struct {
	Controller Service `json:"controller,omitempty"`
	Router     Service `json:"router,omitempty"`
}

type Service struct {
	Type        string            `json:"type,omitempty"`
	Address     string            `json:"address,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type Replicas struct {
	Controller int32 `yaml:"controller"`
}

// Credentials credentials used to log into docker when deploying a local stack
type Credentials struct {
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
}

type Auth struct {
	URL              string `yaml:"url"`
	Realm            string `yaml:"realm"`
	SSL              string `yaml:"ssl"`
	RealmKey         string `yaml:"realmKey"`
	ControllerClient string `yaml:"controllerClient"`
	ControllerSecret string `yaml:"controllerSecret"`
	ViewerClient     string `yaml:"viewerClient"`
}

type Database struct {
	Provider     string  `yaml:"provider,omitempty"`
	Host         string  `yaml:"host,omitempty"`
	Port         int     `yaml:"port,omitempty"`
	User         string  `yaml:"user,omitempty"`
	Password     string  `yaml:"password,omitempty"`
	DatabaseName string  `yaml:"databaseName,omitempty"`
	SSL          *bool   `yaml:"ssl,omitempty"`
	CA           *string `yaml:"ca,omitempty"`
}

type Registry struct {
	URL          *string `yaml:"url"`
	Private      *bool   `yaml:"private"`
	Username     *string `yaml:"username"`
	Password     *string `yaml:"password"`
	Email        *string `yaml:"email"`
	RequiresCert *bool   `yaml:"requiresCert"`
	Certificate  *string `yaml:"certificate,omitempty"`
	ID           int     `yaml:"id"`
}

type Volume struct {
	Name        string   `json:"name" yaml:"name"`
	Agents      []string `json:"agents" yaml:"agents"`
	Source      string   `json:"source" yaml:"source"`
	Destination string   `json:"destination" yaml:"destination"`
	Permissions string   `json:"permissions" yaml:"permissions"`
}

// AgentConfiguration contains configuration information for a deployed agent
type AgentConfiguration struct {
	Name                      string  `json:"name,omitempty" yaml:"name"`
	Location                  string  `json:"location,omitempty" yaml:"location"`
	Latitude                  float64 `json:"latitude,omitempty" yaml:"latitude"`
	Longitude                 float64 `json:"longitude,omitempty" yaml:"longitude"`
	Description               string  `json:"description,omitempty" yaml:"description"`
	FogType                   *string `json:"fogType,omitempty" yaml:"agentType"`
	client.AgentConfiguration `yaml:",inline"`
}

type AgentInfo struct {
	AgentConfiguration
	AgentStatus
}

type AgentStatus struct {
	LastActive            int64   `json:"lastActive" yaml:"lastActive"`
	DaemonStatus          string  `json:"daemonStatus" yaml:"daemonStatus"`
	SecurityStatus        string  `json:"securityStatus" yaml:"securityStatus"`
	SecurityViolationInfo string  `json:"securityViolationInfo" yaml:"securityViolationInfo"`
	WarningMessage        string  `json:"warningMessage" yaml:"warningMessage"`
	UptimeMs              int64   `json:"daemonOperatingDuration" yaml:"uptime"`
	MemoryUsage           float64 `json:"memoryUsage" yaml:"memoryUsage"`
	DiskUsage             float64 `json:"diskUsage" yaml:"diskUsage"`
	CPUUsage              float64 `json:"cpuUsage" yaml:"cpuUsage"`
	SystemAvailableMemory float64 `json:"systemAvailableMemory" yaml:"systemAvailableMemory"`
	SystemAvailableDisk   float64 `json:"systemAvailableDisk" yaml:"systemAvailableDisk"`
	SystemTotalCPU        float64 `json:"systemTotalCPU" yaml:"systemTotalCPU"`
	MemoryViolation       string  `json:"memoryViolation" yaml:"memoryViolation"`
	DiskViolation         string  `json:"diskViolation" yaml:"diskViolation"`
	CPUViolation          string  `json:"cpuViolation" yaml:"cpuViolation"`
	RepositoryStatus      string  `json:"repositoryStatus" yaml:"repositoryStatus"`
	LastStatusTimeMsUTC   int64   `json:"lastStatusTime" yaml:"lastStatusTime"`
	IPAddress             string  `json:"ipAddress" yaml:"ipAddress"`
	IPAddressExternal     string  `json:"ipAddressExternal" yaml:"ipAddressExternal"`
	ProcessedMessaged     int64   `json:"processedMessages" yaml:"ProcessedMessages"`
	MessageSpeed          float64 `json:"messageSpeed" yaml:"messageSpeed"`
	LastCommandTimeMsUTC  int64   `json:"lastCommandTime" yaml:"lastCommandTime"`
	Version               string  `json:"version" yaml:"version"`
	IsReadyToUpgrade      bool    `json:"isReadyToUpgrade" yaml:"isReadyToUpgrade"`
	IsReadyToRollback     bool    `json:"isReadyToRollback" yaml:"isReadyToRollback"`
	Tunnel                string  `json:"tunnel" yaml:"tunnel"`
	VolumeMounts          []VolumeMount
	GpsStatus             string `json:"gpsStatus" yaml:"gpsStatus"`
}

type EdgeResource struct {
	Name              string
	Version           string                     `yaml:"version"`
	Description       string                     `yaml:"description"`
	InterfaceProtocol string                     `yaml:"interfaceProtocol"`
	Interface         *EdgeResourceHTTPInterface `yaml:"interface,omitempty"` // TODO: Make this generic to support multiple interfaces protocols
	Display           *Display                   `yaml:"display,omitempty"`
	OrchestrationTags []string                   `yaml:"orchestrationTags"`
	Custom            map[string]interface{}     `yaml:"custom"`
}

type EdgeResourceHTTPInterface = client.HTTPEdgeResource

type Display = client.EdgeResourceDisplay
type HTTPEndpoint = client.HTTPEndpoint

// FogTypeStringMap map human readable fog type to Controller fog type
var FogTypeStringMap = map[string]int64{
	"auto": 0,
	"x86":  1,
	"arm":  2,
}

// FogTypeIntMap map Controller fog type to human readable fog type
var FogTypeIntMap = map[int]string{
	0: "auto",
	1: "x86",
	2: "arm",
}

type K8SControllerConfig struct {
	PidBaseDir    string `yaml:"pidBaseDir,omitempty"`
	EcnViewerPort int    `yaml:"ecnViewerPort,omitempty"`
	EcnViewerURL  string `yaml:"ecnViewerUrl,omitempty"`
	LogLevel      string `yaml:"logLevel,omitempty"`
	Https         *bool  `yaml:"https,omitempty"`
	SecretName    string `yaml:"secretName,omitempty"`
}

type RemoteControllerConfig struct {
	PidBaseDir    string           `yaml:"pidBaseDir,omitempty"`
	EcnViewerPort int              `yaml:"ecnViewerPort,omitempty"`
	EcnViewerURL  string           `yaml:"ecnViewerUrl,omitempty"`
	LogLevel      string           `yaml:"logLevel,omitempty"`
	Https         *Https           `yaml:"https,omitempty"`
	SiteCA        *SiteCertificate `yaml:"siteCA,omitempty"`  // router site CA
	LocalCA       *SiteCertificate `yaml:"localCA,omitempty"` // router local CA
}

type Https struct {
	Enabled *bool  `yaml:"enabled,omitempty"`
	CACert  string `yaml:"caCert,omitempty"`  // base64 encoded
	TLSCert string `yaml:"tlsCert,omitempty"` // base64 encoded
	TLSKey  string `yaml:"tlsKey,omitempty"`  // base64 encoded
}

type SiteCertificate struct {
	TLSCert string `yaml:"tlsCert,omitempty"` // base64 encoded
	TLSKey  string `yaml:"tlsKey,omitempty"`  // base64 encoded
}

type RouterIngress struct {
	Address      string `yaml:"address,omitempty"`
	MessagePort  int    `yaml:"messagePort,omitempty"`
	InteriorPort int    `yaml:"interiorPort,omitempty"`
	EdgePort     int    `yaml:"edgePort,omitempty"`
}

type ControllerIngress struct {
	Annotations      map[string]string `yaml:"annotations,omitempty"`
	IngressClassName string            `yaml:"ingressClassName,omitempty"`
	Host             string            `yaml:"host,omitempty"`
	SecretName       string            `yaml:"secretName,omitempty"`
}

type Ingress struct {
	Address string `yaml:"address,omitempty"`
}

type Ingresses struct {
	Controller ControllerIngress `yaml:"controller,omitempty"`
	Router     RouterIngress     `yaml:"router,omitempty"`
}

// type RouterConfig struct {
// }

type Secret struct {
	Name string            `yaml:"name,omitempty"`
	Type string            `yaml:"type,omitempty"`
	Data map[string]string `yaml:"data,omitempty"`
}

type ConfigMap struct {
	Name      string            `yaml:"name,omitempty"`
	Immutable bool              `yaml:"immutable"`
	Data      map[string]string `yaml:"data,omitempty"`
}

type ClusterService struct {
	Name            string   `yaml:"name,omitempty"`
	Tags            []string `yaml:"tags,omitempty"`
	Type            string   `yaml:"type,omitempty"`
	Resource        string   `yaml:"resource,omitempty"`
	TargetPort      int      `yaml:"targetPort,omitempty"`
	BridgePort      int      `yaml:"bridgePort,omitempty"`
	DefaultBridge   string   `yaml:"defaultBridge,omitempty"`
	K8sType         string   `yaml:"k8sType,omitempty"`
	ServiceEndpoint string   `yaml:"serviceEndpoint,omitempty"`
	ServicePort     int      `yaml:"servicePort,omitempty"`
}

type VolumeMount struct {
	Name          string `yaml:"name,omitempty"`
	UUID          string `yaml:"uuid,omitempty"`
	ConfigMapName string `yaml:"configMapName,omitempty"`
	SecretName    string `yaml:"secretName,omitempty"`
	Version       int    `yaml:"version,omitempty"`
}

type CertificateInfo struct {
	// Name             string                 `json:"name"`
	Subject      string    `json:"subject" yaml:"subject"`
	Hosts        string    `json:"hosts" yaml:"hosts"`
	IsCA         bool      `json:"isCA" yaml:"isCA"`
	ValidFrom    time.Time `json:"validFrom" yaml:"validFrom"`
	ValidTo      time.Time `json:"validTo" yaml:"validTo"`
	SerialNumber string    `json:"serialNumber" yaml:"serialNumber"`
	CAName       *string   `json:"caName" yaml:"caName"`
	// CertificateChain []CertificateChainItem `json:"certificateChain" yaml:"certificateChain"`
	DaysRemaining int    `json:"daysRemaining" yaml:"daysRemaining"`
	IsExpired     bool   `json:"isExpired" yaml:"isExpired"`
	Certificate   string `json:"certificate" yaml:"certificate"`
	PrivateKey    string `json:"privateKey" yaml:"privateKey"`
}

type CertificateChainItem struct {
	Name    string `json:"name" yaml:"name"`
	Subject string `json:"subject" yaml:"subject"`
}

type CAInfo struct {
	// Name         string          `json:"name"`
	Subject      string    `json:"subject" yaml:"subject"`
	IsCA         bool      `json:"isCA" yaml:"isCA"`
	ValidFrom    time.Time `json:"validFrom" yaml:"validFrom"`
	ValidTo      time.Time `json:"validTo" yaml:"validTo"`
	SerialNumber string    `json:"serialNumber" yaml:"serialNumber"`
	Certificate  string    `json:"certificate" yaml:"certificate"`
	PrivateKey   string    `json:"privateKey" yaml:"privateKey"`
}

type CertificateData struct {
	Certificate string `json:"certificate" yaml:"certificate"`
	PrivateKey  string `json:"privateKey" yaml:"privateKey"`
}

// Certificate Types
type CertificateCreateRequest struct {
	Name       string              `json:"name" yaml:"name"`
	Subject    string              `json:"subject" yaml:"subject"`
	Hosts      string              `json:"hosts" yaml:"hosts"`
	Expiration int                 `json:"expiration,omitempty" yaml:"expiration,omitempty"`
	CA         CertificateCreateCA `json:"ca" yaml:"ca"`
}

type CertificateCreateResponse struct {
	Name      string    `json:"name" yaml:"name"`
	Subject   string    `json:"subject" yaml:"subject"`
	Hosts     string    `json:"hosts" yaml:"hosts"`
	ValidFrom time.Time `json:"validFrom" yaml:"validFrom"`
	ValidTo   time.Time `json:"validTo" yaml:"validTo"`
	CAName    string    `json:"caName" yaml:"caName"`
}

type CertificateCACreateResponse struct {
	Name      string    `json:"name" yaml:"name"`
	Subject   string    `json:"subject" yaml:"subject"`
	Type      string    `json:"type" yaml:"type"`
	ValidFrom time.Time `json:"validFrom" yaml:"validFrom"`
	ValidTo   time.Time `json:"validTo" yaml:"validTo"`
}

type CertificateCreateCA struct {
	Type       string `json:"type" yaml:"type"`
	SecretName string `json:"secretName,omitempty" yaml:"secretName,omitempty"`
}

type CACreateRequest struct {
	Name       string `json:"name" yaml:"name"`
	Subject    string `json:"subject,omitempty" yaml:"subject,omitempty"`
	Expiration int    `json:"expiration,omitempty" yaml:"expiration,omitempty"`
	Type       string `json:"type" yaml:"type" yaml:"type"`
	SecretName string `json:"secretName,omitempty" yaml:"secretName,omitempty"`
}
