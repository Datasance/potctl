package deployairgap

import (
	"fmt"

	"github.com/datasance/iofog-go-sdk/v3/pkg/client"
	"github.com/datasance/potctl/internal/config"
	rsc "github.com/datasance/potctl/internal/resource"
	clientutil "github.com/datasance/potctl/internal/util/client"
	"github.com/datasance/potctl/pkg/util"
)

// RequiredImages represents all images needed for airgap deployment
type RequiredImages struct {
	Controller string
	Agent      string
	RouterX86  string
	RouterARM  string
	Debugger   string
}

// CollectControllerImages collects required images for controller deployment
// For initial deployment: uses YAML/defaults
// For existing control plane: fetches router/debugger from controller catalog items
func CollectControllerImages(namespace string, controlPlane *rsc.RemoteControlPlane, isInitialDeployment bool) (*RequiredImages, error) {
	images := &RequiredImages{}

	// Controller image
	if controlPlane.Package.Container.Image != "" {
		images.Controller = controlPlane.Package.Container.Image
	} else {
		images.Controller = util.GetControllerImage()
	}

	// Router and debugger images
	if isInitialDeployment {
		// Use YAML or defaults
		if controlPlane.SystemMicroservices.Router.X86 != "" {
			images.RouterX86 = controlPlane.SystemMicroservices.Router.X86
		} else {
			images.RouterX86 = util.GetRouterImage()
		}

		if controlPlane.SystemMicroservices.Router.ARM != "" {
			images.RouterARM = controlPlane.SystemMicroservices.Router.ARM
		} else {
			images.RouterARM = util.GetRouterARMImage()
		}

		// Debugger image
		// TODO: Add debugger image to RemoteSystemMicroservices
		images.Debugger = util.GetDebuggerImage()
	} else {
		// Fetch from controller catalog items
		clt, err := clientutil.NewControllerClient(namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to controller: %w", err)
		}

		// Get router catalog item
		routerItem, err := clt.GetCatalogItemByName("router")
		if err != nil {
			return nil, fmt.Errorf("failed to get router catalog item from controller: %w", err)
		}

		// Extract router images from catalog item
		for _, img := range routerItem.Images {
			switch client.AgentTypeIDAgentTypeDict[img.AgentTypeID] {
			case "x86":
				images.RouterX86 = img.ContainerImage
			case "arm":
				images.RouterARM = img.ContainerImage
			}
		}

		// Get debugger catalog item
		debuggerItem, err := clt.GetCatalogItemByName("debugger")
		if err != nil {
			// Debugger is optional, log but don't fail
			util.PrintNotify("Warning: Could not fetch debugger catalog item from controller. Debugger image will not be transferred.")
			images.Debugger = ""
		} else {
			// Extract debugger image (typically x86, but check both)
			for _, img := range debuggerItem.Images {
				if client.AgentTypeIDAgentTypeDict[img.AgentTypeID] == "x86" {
					images.Debugger = img.ContainerImage
					break
				} else if client.AgentTypeIDAgentTypeDict[img.AgentTypeID] == "arm" {
					images.Debugger = img.ContainerImage
					break
				}
			}
		}
	}

	return images, nil
}

// CollectAgentImages collects required images for agent deployment
// For initial deployment: uses YAML/defaults (only for RemoteControlPlane)
// For existing control plane: fetches router/debugger from controller catalog items
// controlPlane can be nil for Kubernetes or other non-remote control planes
func CollectAgentImages(namespace string, agent *rsc.RemoteAgent, controlPlane *rsc.RemoteControlPlane, isInitialDeployment bool) (*RequiredImages, error) {
	images := &RequiredImages{}

	// Agent image
	if agent.Package.Container.Image != "" {
		images.Agent = agent.Package.Container.Image
	} else {
		images.Agent = util.GetAgentImage()
	}

	// Router and debugger images
	if isInitialDeployment && controlPlane != nil {
		// Use YAML or defaults (only for RemoteControlPlane initial deployment)
		if controlPlane.SystemMicroservices.Router.X86 != "" {
			images.RouterX86 = controlPlane.SystemMicroservices.Router.X86
		} else {
			images.RouterX86 = util.GetRouterImage()
		}

		if controlPlane.SystemMicroservices.Router.ARM != "" {
			images.RouterARM = controlPlane.SystemMicroservices.Router.ARM
		} else {
			images.RouterARM = util.GetRouterARMImage()
		}

		// Debugger image - no default in YAML structure
		images.Debugger = ""
	} else {
		// Fetch from controller catalog items (for existing control plane or non-remote control planes)
		clt, err := clientutil.NewControllerClient(namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to controller: %w", err)
		}

		// Get router catalog item
		routerItem, err := clt.GetCatalogItemByName("Router")
		if err != nil {
			return nil, fmt.Errorf("failed to get router catalog item from controller: %w", err)
		}

		// Extract router images from catalog item
		for _, img := range routerItem.Images {
			switch client.AgentTypeIDAgentTypeDict[img.AgentTypeID] {
			case "x86":
				images.RouterX86 = img.ContainerImage
			case "arm":
				images.RouterARM = img.ContainerImage
			}
		}

		// Get debugger catalog item
		debuggerItem, err := clt.GetCatalogItemByName("Debug")
		if err != nil {
			// Debugger is optional, log but don't fail
			util.PrintNotify("Warning: Could not fetch debugger catalog item from controller. Debugger image will not be transferred.")
			images.Debugger = ""
		} else {
			// Extract debugger image (typically x86, but check both)
			for _, img := range debuggerItem.Images {
				if client.AgentTypeIDAgentTypeDict[img.AgentTypeID] == "x86" {
					images.Debugger = img.ContainerImage
					break
				} else if client.AgentTypeIDAgentTypeDict[img.AgentTypeID] == "arm" {
					images.Debugger = img.ContainerImage
					break
				}
			}
		}
	}

	return images, nil
}

// GetImageForPlatform returns the appropriate image based on platform
func GetImageForPlatform(images *RequiredImages, platform string) (string, error) {
	switch platform {
	case PlatformAMD64:
		return images.RouterX86, nil
	case PlatformARM64:
		return images.RouterARM, nil
	default:
		return "", util.NewInputError(fmt.Sprintf("unsupported platform %s", platform))
	}
}

// IsInitialDeployment checks if this is an initial control plane deployment
func IsInitialDeployment(namespace string) (bool, error) {
	ns, err := config.GetNamespace(namespace)
	if err != nil {
		return false, err
	}
	// If control plane exists and has controllers, this is not initial deployment
	controlPlane, err := ns.GetControlPlane()
	if err != nil {
		// No control plane exists, this is initial deployment
		return true, nil
	}
	// If controllers exist, this is not initial deployment
	controllers := controlPlane.GetControllers()
	return len(controllers) == 0, nil
}
