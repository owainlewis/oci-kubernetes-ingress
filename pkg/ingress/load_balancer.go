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
	details := loadbalancer.CreateLoadBalancerDetails{
		CompartmentId: common.String(specification.GetLoadBalancerCompartment()),
		DisplayName:   common.String(specification.Name),
		ShapeName:     common.String(specification.GetLoadBalancerShape()),
		IsPrivate:     common.Bool(false),
		SubnetIds:     specification.GetLoadBalancerSubnets(),
		//BackendSets:   spec.BackendSets,
		//Listeners:     spec.Listeners,
		//Certificates:  certs,
	}

	request := loadbalancer.CreateLoadBalancerRequest{
		CreateLoadBalancerDetails: details,
	}

	ctx := context.Background()
	return svc.client.CreateLoadBalancer(ctx, request)
}

// DeleteLoadBalancer will delete a load balancer from OCI.
func (svc *LoadBalancerService) DeleteLoadBalancer(name string) error {
	//request := loadbalancer.DeleteLoadBalancerRequest{
	//LoadBalancerId: common.String(id),
	//}

	//ctx := context.Background()
	// _, err := svc.client.DeleteLoadBalancer(ctx, request)
	// return err

	return nil
}

// GetLoadBalancer will find an OCI load balancer based on a specification.
func (svc *LoadBalancerService) GetLoadBalancer(specification Specification) (loadbalancer.LoadBalancer, error) {
	request := loadbalancer.ListLoadBalancersRequest{
		CompartmentId: common.String(specification.GetLoadBalancerCompartment()),
		DisplayName:   common.String(specification.Name),
	}

	ctx := context.Background()
	response, err := svc.client.ListLoadBalancers(ctx, request)
	if err != nil {
		return loadbalancer.LoadBalancer{}, err
	}

	for _, lb := range response.Items {
		if common.PointerString(lb.DisplayName) == specification.Name {
			return lb, nil
		}
	}

	return loadbalancer.LoadBalancer{}, fmt.Errorf("Could not find load balancer with display name %s", specification.Name)
}
