package annotations

import (
	extensions "k8s.io/api/extensions/v1beta1"
)

const (
	annotationKubernetesIngressClass = "kubernetes.io/ingress.class"
	defaultIngressClass              = "oci"
	// LoadBalancerShape is an annotation for
	// specifying the Shape of a load balancer. The shape is a template that
	// determines the load balancer's total pre-provisioned maximum capacity
	// (bandwidth) for ingress plus egress traffic. Available shapes include
	// "100Mbps", "400Mbps", and "8000Mbps".
	LoadBalancerShape = "oci.ingress.kubernetes.io/loadbalancer-shape"
	// LoadBalancerSubnets defines a comma separated list of subnets for the load balancer
	LoadBalancerSubnets = "oci.ingress.kubernetes.io/loadbalancer-subnets"
)

// HasOCIIngressAnnotation returns true if an ingress object has the OCI ingress class
func HasOCIIngressAnnotation(ingress *extensions.Ingress) bool {
	actualIngressClass := ingress.GetAnnotations()[annotationKubernetesIngressClass]
	return actualIngressClass == defaultIngressClass
}
