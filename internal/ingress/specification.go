package ingress

import (
	"k8s.io/api/extensions/v1beta1"
)

// A specification describes the desired state of the OCI load balancer
type Specification struct {
	Name     string
	Shape    string
	Subnets  []string
	Internal bool
}

func NewSpecificationFromIngress(ingress *v1beta1.Ingress) Specification {
	return Specification{}
}
