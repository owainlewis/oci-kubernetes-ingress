package common

import (
	"fmt"
)

// GetLoadBalancerName returns the name of a load balancer which is
// derived by combinging the ingress namespace and ingress name
func GetLoadBalancerName(namespace string, ingressName string) string {
	name := fmt.Sprintf("%s-%s", namespace, ingressName)
	if len(name) > 1024 {
		// 1024 is the max length for display name
		// https://docs.us-phoenix-1.oraclecloud.com/api/#/en/loadbalancer/20170115/requests/UpdateLoadBalancerDetails
		name = name[:1024]
	}

	return name
}
