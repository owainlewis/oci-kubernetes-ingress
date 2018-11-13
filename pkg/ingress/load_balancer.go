package ingress

import (
	"context"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/loadbalancer"
)

// LoadBalancerService is responsible for managing OCI Load Balancers.
type LoadBalancerService struct {
	client loadbalancer.LoadBalancerClient
}

// NewLoadBalancerService will create a service to manage OCI Load Balancers.
func NewLoadBalancerService() (*LoadBalancerService, error) {
	client, err := loadbalancer.NewLoadBalancerClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		return nil, err
	}

	return &LoadBalancerService{client: client}, nil
}

// CreateLoadBalancer will create a new OCI load balancer to handle ingress traffic.
func (svc *LoadBalancerService) CreateLoadBalancer(specification Specification) (loadbalancer.CreateLoadBalancerResponse, error) {
	details := loadbalancer.CreateLoadBalancerDetails{
		CompartmentId: common.String("ocid1.compartment.oc1..aaaaaaaaob4ckouj3cjmf36ifjkff33wvln5fnnarumafqzpqq7tmbig2n5q"),
		DisplayName:   common.String(specification.Name),
		ShapeName:     common.String(specification.GetLoadBalancerShape()),
		IsPrivate:     common.Bool(false),
		SubnetIds: []string{
			"ocid1.subnet.oc1.uk-london-1.aaaaaaaaqalydfvmgw7pdw3tittizpoyondib7hedwayyswrrfcrsmc4j7dq",
			"ocid1.subnet.oc1.uk-london-1.aaaaaaaa2tqtopdpynhbjglh3szj2j6h6pwwwohrcanbeyj6dpbiboyuvrza",
		},
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

func (svc *LoadBalancerService) getLoadBalancerByName(name string) {

}
