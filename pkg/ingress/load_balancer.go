package ingress

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/loadbalancer"
	"github.com/owainlewis/oci-kubernetes-ingress/pkg/config"
)

// LoadBalancerService is responsible for managing OCI Load Balancers.
type LoadBalancerService struct {
	client loadbalancer.LoadBalancerClient
	config config.Config
}

// NewLoadBalancerService will create a service to manage OCI Load Balancers.
func NewLoadBalancerService(conf config.Config) (*LoadBalancerService, error) {
	client, err := loadbalancer.NewLoadBalancerClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		return nil, err
	}

	return &LoadBalancerService{client: client, config: conf}, nil
}

// CreateLoadBalancer will create a new OCI load balancer to handle ingress traffic.
func (svc *LoadBalancerService) CreateLoadBalancer(specification Specification) (loadbalancer.CreateLoadBalancerResponse, error) {
	request := loadbalancer.CreateLoadBalancerRequest{
		CreateLoadBalancerDetails: loadbalancer.CreateLoadBalancerDetails{
			CompartmentId: common.String(specification.GetLoadBalancerCompartment()),
			DisplayName:   common.String(GetLoadBalancerUniqueName(specification.Ingress)),
			ShapeName:     common.String(specification.GetLoadBalancerShape()),
			IsPrivate:     common.Bool(false),
			SubnetIds:     specification.GetLoadBalancerSubnets(),

			//BackendSets:   spec.BackendSets,
			//Listeners:     spec.Listeners,
			//Certificates:  certs,
		},
	}

	ctx := context.Background()
	return svc.client.CreateLoadBalancer(ctx, request)
}

// DeleteLoadBalancer will delete a load balancer from OCI.
func (svc *LoadBalancerService) DeleteLoadBalancer(compartment, name string) error {
	lb, err := svc.GetLoadBalancer(compartment, name)
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = svc.client.DeleteLoadBalancer(ctx, loadbalancer.DeleteLoadBalancerRequest{
		LoadBalancerId: lb.Id,
	})
	return err
}

// GetLoadBalancer will find an OCI load balancer based on a specification.
// TODO we are relying on display name which is not unique in OCI
func (svc *LoadBalancerService) GetLoadBalancer(compartment string, name string) (loadbalancer.LoadBalancer, error) {
	request := loadbalancer.ListLoadBalancersRequest{
		CompartmentId: common.String(compartment),
		DisplayName:   common.String(name),
	}

	ctx := context.Background()
	response, err := svc.client.ListLoadBalancers(ctx, request)
	if err != nil {
		return loadbalancer.LoadBalancer{}, err
	}

	for _, lb := range response.Items {
		if lb.DisplayName != nil && *lb.DisplayName == name {
			return lb, nil
		}
	}

	return loadbalancer.LoadBalancer{}, fmt.Errorf("Could not find load balancer with display name %s", name)
}
