package ingress

func (spec Specification) withIngressAnnotations(k, v string) Specification {
	annotations := map[string]string{}
	annotations[k] = v

	spec.Ingress.Annotations = annotations
	return spec
}
