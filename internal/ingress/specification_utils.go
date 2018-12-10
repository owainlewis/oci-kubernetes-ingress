package ingress

func (spec Specification) withIngressAnnotation(k, v string) Specification {
	annotations := spec.Ingress.Annotations
	annotations[k] = v

	spec.Ingress.Annotations = annotations
	return spec
}
