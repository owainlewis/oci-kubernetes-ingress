package main

import (
	"context"
	"time"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/loadbalancer"
)

type loadBalancerService struct {
	client loadbalancer.LoadBalancerClient
}

// NewLoadBalancerService ...
func NewLoadBalancerService() (*loadBalancerService, error) {
	client, err := loadbalancer.NewLoadBalancerClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		return nil, err
	}

	return &loadBalancerService{client: client}, nil
}

func (svc *loadBalancerService) Create() {
	details := loadbalancer.CreateLoadBalancerDetails{
		CompartmentId: common.String("ocid1.compartment.oc1..aaaaaaaaob4ckouj3cjmf36ifjkff33wvln5fnnarumafqzpqq7tmbig2n5q"),
		DisplayName:   common.String("mylb"),
		ShapeName:     common.String("100Mbps"),
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
	svc.client.CreateLoadBalancer(ctx, request)
}

func (svc *loadBalancerService) Delete() {

}

func (svc *loadBalancerService) Update() {

}

func main() {

	svc, _ := NewLoadBalancerService()

	svc.Create()

	time.Sleep(2 * time.Minute)
}
