package connectcontrolplane

import (
	"github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client"
	rsc "github.com/datasance/potctl/internal/resource"
	"github.com/datasance/potctl/pkg/util"
)

func Connect(ctrlPlane rsc.ControlPlane, endpoint string, ns *rsc.Namespace) error {
	// Connect to Controller
	baseURL, err := util.GetBaseURL(endpoint)
	if err != nil {
		return err
	}
	ctrl, err := client.NewAndLogin(client.Options{BaseURL: baseURL}, ctrlPlane.GetUser().Email, ctrlPlane.GetUser().GetRawPassword())
	if err != nil {
		return err
	}

	// Get Agents
	listAgentsResponse, err := ctrl.ListAgents(client.ListAgentsRequest{})
	if err != nil {
		return err
	}

	// Update Agents config
	for idx := range listAgentsResponse.Agents {
		agent := &listAgentsResponse.Agents[idx]
		agentConfig := rsc.RemoteAgent{
			Name: agent.Name,
			UUID: agent.UUID,
			Host: agent.IPAddressExternal,
		}
		if err := ns.AddAgent(&agentConfig); err != nil {
			return err
		}
	}
	return nil
}
