package ingress

import (
	"testing"

	"github.com/owainlewis/oci-kubernetes-ingress/internal/annotations"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/config"
	v1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func newIngressBackend(svcName string, svcPort int) *extensions.IngressBackend {
	return &extensions.IngressBackend{
		ServiceName: svcName,
		ServicePort: intstr.FromInt(svcPort),
	}
}

func newIngress(defaultBackend *extensions.IngressBackend) *extensions.Ingress {
	return &extensions.Ingress{
		Spec: extensions.IngressSpec{
			Backend: defaultBackend,
		},
	}
}

func newConfiguration() config.Config {
	return config.Config{
		Loadbalancer: config.LoadbalancerConfig{
			Compartment: "default",
			Subnets: []string{
				"ocid.subnet.a",
				"ocid.subnet.b",
			},
		},
	}
}

func TestGetLoadBalancerShapeCustomAnnotations(t *testing.T) {
	configuration := newConfiguration()
	defaultBackend := newIngressBackend("nginx", 80)
	ingress := newIngress(defaultBackend)

	nodes := []*v1.Node{}

	specification := NewSpecification(configuration, ingress, nodes).
		withIngressAnnotation(annotations.LoadBalancerShape, "400mbps")

	shape := specification.GetLoadBalancerShape()
	if shape != "400mbps" {
		t.Errorf("Expected LB shape to be 400mbps but got: %s", shape)
	}
}
