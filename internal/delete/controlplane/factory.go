package deletecontrolplane

import (
	"github.com/datasance/potctl/internal/config"
	deletek8scontrolplane "github.com/datasance/potctl/internal/delete/controlplane/k8s"
	deletelocalcontrolplane "github.com/datasance/potctl/internal/delete/controlplane/local"
	deleteremotecontrolplane "github.com/datasance/potctl/internal/delete/controlplane/remote"
	"github.com/datasance/potctl/internal/execute"
	rsc "github.com/datasance/potctl/internal/resource"
	"github.com/datasance/potctl/pkg/util"
)

func NewExecutor(namespace string) (execute.Executor, error) {
	ns, err := config.GetNamespace(namespace)
	if err != nil {
		return nil, err
	}
	baseControlPlane, err := ns.GetControlPlane()
	if err != nil {
		return nil, err
	}

	// //Deactivation for entitlement will be done

	// user := baseControlPlane.GetUser()
	// util.DeactivateEntitlementDatasance(user.SubscriptionKey, namespace, user.Email)

	switch baseControlPlane.(type) {
	case *rsc.KubernetesControlPlane:
		return deletek8scontrolplane.NewExecutor(namespace)
	case *rsc.RemoteControlPlane:
		return deleteremotecontrolplane.NewExecutor(namespace)
	case *rsc.LocalControlPlane:
		return deletelocalcontrolplane.NewExecutor(namespace)
	}
	return nil, util.NewError("Could not convert Control Plane to dynamic type")
}
