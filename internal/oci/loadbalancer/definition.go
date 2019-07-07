package loadbalancer

import (
	extensions "k8s.io/api/extensions/v1beta1"

	"github.com/owainlewis/oci-kubernetes-ingress/internal/common"
)

// Definition is a structure representing the definition
// of a loadbalancer
type Definition struct {
	Name  string
	Shape string
}

// NewLoadBalancerDefinition builds a load balancer defintion which is an internalised
// simplified representation of a load balancer derived by combinging together all
// the relevant information needed to create an OCI load balancer for a particular ingress.
func NewLoadBalancerDefinition(ingress *extensions.Ingress) Definition {
	lbName := common.GetLoadBalancerName(ingress.Namespace, ingress.Name)
	lbShape := getLoadBalancerShape(ingress)

	return Definition{
		Name:  lbName,
		Shape: lbShape,
	}
}

func getLoadBalancerShape(ingress *extensions.Ingress) string {
	return "100Mbps"
}
