package ingress

import (
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
)

// Specification describes the desired state of the OCI load balancer
type Specification struct {
	Name    string
	Ingress *v1beta1.Ingress
	Nodes   []*core_v1.Node
}

// NewSpecification creates a new load balancer specification for a
// given Ingress
func NewSpecification(name string, ingress *v1beta1.Ingress, nodes []*core_v1.Node) Specification {
	return Specification{Name: name, Ingress: ingress, Nodes: nodes}
}

// GetLoadBalancerShape will return the load balancer shape required.
// The shape can be controlled by setting ingress object annotations.
func (spec Specification) GetLoadBalancerShape() string {
	return "100Mbps"
}
