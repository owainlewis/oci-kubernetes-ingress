package loadbalancer

type LoadBalancerDefinition struct {
	Name  string
	Shape string
}

func NewLoadBalancerDefinition(name string, shape string) LoadBalancerDefinition {
	return LoadBalancerDefinition{
		Name:  name,
		Shape: shape,
	}
}
