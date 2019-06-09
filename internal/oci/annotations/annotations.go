package annotations

const (
	// LoadBalancerVisibility is an annotation for
	// specifying that a load balancer should be public or private.
	// By default all load balancers will be public.
	// Values an be one of "public" or "private"
	LoadBalancerVisibility = "ingress.beta.kubernetes.io/oci-load-balancer-visibility"
	// LoadBalancerShape is an annotation for
	// specifying the shape of a load balancer. Available shapes include
	// "100Mbps", "400Mbps", and "8000Mbps".
	LoadBalancerShape = "ingress.beta.kubernetes.io/oci-load-balancer-shape"
	// LoadBalancerCompartment allows for load balancers to be created in a compartment
	// different to that specified in config.
	LoadBalancerCompartment = "ingress.beta.kubernetes.io/oci-load-balancer-compartment"
)
