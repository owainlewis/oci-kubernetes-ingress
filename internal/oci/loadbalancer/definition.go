package loadbalancer

import (
	"fmt"

	extensions "k8s.io/api/extensions/v1beta1"
)

type LoadBalancerDefinition struct {
	Name  string
	Shape string
}

func NewLoadBalancerDefinition(ingress *extensions.Ingress) LoadBalancerDefinition {
	lbName := getLoadBalancerName(ingress)
	lbShape := getLoadBalancerShape(ingress)

	return LoadBalancerDefinition{
		Name:  lbName,
		Shape: lbShape,
	}
}

func getLoadBalancerShape(ingress *extensions.Ingress) string {
	return "100Mbps"
}

func getLoadBalancerName(ingress *extensions.Ingress) string {
	name := fmt.Sprintf("%s-%s", ingress.Namespace, ingress.Name)
	if len(name) > 1024 {
		// 1024 is the max length for display name
		// https://docs.us-phoenix-1.oraclecloud.com/api/#/en/loadbalancer/20170115/requests/UpdateLoadBalancerDetails
		name = name[:1024]
	}

	return name
}
