package annotations

import (
	extensions "k8s.io/api/extensions/v1beta1"
)

const (
	annotationKubernetesIngressClass = "kubernetes.io/ingress.class"
	defaultIngressClass              = "oci"
)

// HasOCIIngressAnnotation returns true if an ingress object has the OCI ingress class
func HasOCIIngressAnnotation(ingress *extensions.Ingress) bool {
	actualIngressClass := ingress.GetAnnotations()[annotationKubernetesIngressClass]
	return actualIngressClass == defaultIngressClass
}
