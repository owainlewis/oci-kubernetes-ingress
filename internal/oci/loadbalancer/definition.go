package loadbalancer

import (
	"github.com/oracle/oci-go-sdk/loadbalancer"
	util "github.com/owainlewis/oci-kubernetes-ingress/internal/common"
	"k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
)

// Definition is a structure representing the definition
// of a loadbalancer
type Definition struct {
	Name        string
	Shape       string
	Subnets     []string
	Listeners   map[string]loadbalancer.ListenerDetails
	BackendSets map[string]loadbalancer.BackendSetDetails
}

// NewLoadBalancerDefinition builds a load balancer defintion which is an internalised
// simplified representation of a load balancer derived by combinging together all
// the relevant information needed to create an OCI load balancer for a particular ingress.
func NewLoadBalancerDefinition(ingress *extensions.Ingress, nodes []*v1.Node, subnets []string) Definition {
	lbName := util.GetLoadBalancerName(ingress.Namespace, ingress.Name)
	lbShape := getLoadBalancerShape(ingress)

	listeners, err := getListeners(ingress)
	if err != nil {

	}

	backendsets, err := getBackendSets(ingress)
	if err != nil {

	}

	return Definition{
		Name:        lbName,
		Shape:       lbShape,
		Subnets:     subnets,
		Listeners:   listeners,
		BackendSets: backendsets,
	}
}

func getLoadBalancerShape(ingress *extensions.Ingress) string {
	return "100Mbps"
}

func getBackendSets(ingress *extensions.Ingress) (map[string]loadbalancer.BackendSetDetails, error) {
	return nil, nil
}

func getPaths(ingress *extensions.Ingress) {
	// for _, rule := range ingress.Spec.Rules {
	// 	for _, _ := range rule.HTTP.Paths {

	// 	}
	// }
}

func getListeners(ingress *extensions.Ingress) (map[string]loadbalancer.ListenerDetails, error) {

	// for _, rule := range ingress.Spec.Rules {
	// 	for _, _ := range rule.HTTP.Paths {

	// 	}
	// }

	// loadbalancer.ListenerDetails{
	// 	DefaultBackendSetName: common.String("backendsetforlistener"),
	// 	Protocol:              common.String("HTTP"),
	// 	Port:                  common.Int(80),
	// 	SslConfiguration:      nil,
	// }

	return nil, nil
}
