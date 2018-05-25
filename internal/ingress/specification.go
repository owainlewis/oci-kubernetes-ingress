package ingress

import (
	"k8s.io/api/extensions/v1beta1"
)

// Specification describes the desired state of the OCI load balancer
type Specification struct {
	Name    string
	Shape   string
	Subnets []string
}

// NewSpecificationFromIngress creates a new load balancer specification for a
// given Ingress
func NewSpecificationFromIngress(ingress *v1beta1.Ingress) Specification {
	return Specification{}
}

func (spec *Specification) getBackendSets() {
}

func (spec *Specification) getListeners() {
}

func (spec *Specification) getHostNames() {
}
