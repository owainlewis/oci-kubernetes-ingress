package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/loadbalancer"
)

type loadBalancerService struct {
	client loadbalancer.LoadBalancerClient
}

func NewLoadBalancerService() (*loadBalancerService, error) {
	client, err := loadbalancer.NewLoadBalancerClientWithConfigurationProvider(common.DefaultConfigProvider())

	if err != nil {
		return nil, err
	}

	return &loadBalancerService{client: client}, nil
}

func (svc *loadBalancerService) CreateLoadBalancer() (loadbalancer.CreateLoadBalancerResponse, error) {
	ctx := context.Background()

	req := loadbalancer.CreateLoadBalancerRequest{}

	req.DisplayName = common.String("My Load Balancer")

	return svc.client.CreateLoadBalancer(ctx, req)
}
