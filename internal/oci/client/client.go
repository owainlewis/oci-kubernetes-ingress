package client

import (
	"context"
	"errors"

	"github.com/oracle/oci-go-sdk/loadbalancer"

	"github.com/oracle/oci-go-sdk/common"
)

var errNotFound = errors.New("not found")

// OCIClient is an interface representing the available (simplified) OCI operations
type OCIClient interface {
	CreateLoadBalancer(ctx context.Context, request loadbalancer.CreateLoadBalancerRequest) (loadbalancer.CreateLoadBalancerResponse, error)
	GetLoadBalancerByName(ctx context.Context, compartmentID string, name string) (*loadbalancer.LoadBalancer, error)
}

type ociClient struct {
	loadbalancer loadbalancer.LoadBalancerClient
}

func (oci ociClient) CreateLoadBalancer(ctx context.Context, request loadbalancer.CreateLoadBalancerRequest) (loadbalancer.CreateLoadBalancerResponse, error) {
	return oci.loadbalancer.CreateLoadBalancer(ctx, request)
}

func (oci ociClient) GetLoadBalancerByName(ctx context.Context, compartmentID string, name string) (*loadbalancer.LoadBalancer, error) {
	var page *string
	for {
		resp, err := oci.loadbalancer.ListLoadBalancers(ctx, loadbalancer.ListLoadBalancersRequest{
			CompartmentId: common.String(compartmentID),
			DisplayName:   common.String(name),
			Page:          page,
		})

		if err != nil {
			return nil, err
		}
		for _, lb := range resp.Items {
			if *lb.DisplayName == name {
				return &lb, nil
			}
		}
		if page = resp.OpcNextPage; page == nil {
			break
		}
	}

	return nil, errNotFound
}

func NewOCI(provider common.ConfigurationProvider) (OCIClient, error) {
	lbClient, err := loadbalancer.NewLoadBalancerClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	return &ociClient{
		loadbalancer: lbClient,
	}, nil
}
