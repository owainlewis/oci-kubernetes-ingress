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
	ListLoadBalancers(ctx context.Context, compartmentID string) ([]loadbalancer.LoadBalancer, error)
	CreateLoadBalancer(ctx context.Context, request loadbalancer.CreateLoadBalancerRequest) (loadbalancer.CreateLoadBalancerResponse, error)
	DeleteLoadBalancer(ctx context.Context, id string) (string, error)
	DeleteLoadBalancerByName(ctx context.Context, compartmentID string, name string) (string, error)
	GetLoadBalancerByName(ctx context.Context, compartmentID string, name string) (*loadbalancer.LoadBalancer, error)
}

type ociClient struct {
	loadbalancer loadbalancer.LoadBalancerClient
}

func (oci ociClient) ListLoadBalancers(ctx context.Context, compartmentID string) ([]loadbalancer.LoadBalancer, error) {
	var page *string
	var result []loadbalancer.LoadBalancer
	for {
		resp, err := oci.loadbalancer.ListLoadBalancers(ctx, loadbalancer.ListLoadBalancersRequest{
			CompartmentId:  common.String(compartmentID),
			LifecycleState: loadbalancer.LoadBalancerLifecycleStateActive,
			Page:           page,
		})

		if err != nil {
			return nil, err
		}
		for _, lb := range resp.Items {
			result = append(result, lb)
		}

		if page = resp.OpcNextPage; page == nil {
			break
		}
	}

	return result, nil
}

func (oci ociClient) CreateLoadBalancer(ctx context.Context, request loadbalancer.CreateLoadBalancerRequest) (loadbalancer.CreateLoadBalancerResponse, error) {
	return oci.loadbalancer.CreateLoadBalancer(ctx, request)
}

func (oci ociClient) DeleteLoadBalancer(ctx context.Context, id string) (string, error) {
	resp, err := oci.loadbalancer.DeleteLoadBalancer(ctx, loadbalancer.DeleteLoadBalancerRequest{
		LoadBalancerId: common.String(id),
	})

	if err != nil {
		return "", err
	}

	return *resp.OpcWorkRequestId, nil
}

func (oci ociClient) DeleteLoadBalancerByName(ctx context.Context, compartmentID string, name string) (string, error) {
	lb, err := oci.GetLoadBalancerByName(ctx, compartmentID, name)
	if err != nil || lb.Id == nil {
		return "", err
	}

	return oci.DeleteLoadBalancer(ctx, *lb.Id)
}

func (oci ociClient) GetLoadBalancerByName(ctx context.Context, compartmentID string, name string) (*loadbalancer.LoadBalancer, error) {
	var page *string
	for {
		resp, err := oci.loadbalancer.ListLoadBalancers(ctx, loadbalancer.ListLoadBalancersRequest{
			CompartmentId:  common.String(compartmentID),
			DisplayName:    common.String(name),
			LifecycleState: loadbalancer.LoadBalancerLifecycleStateActive,
			Page:           page,
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

// NewOCI builds a new OCI client
func NewOCI(provider common.ConfigurationProvider) (OCIClient, error) {
	lbClient, err := loadbalancer.NewLoadBalancerClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	return &ociClient{
		loadbalancer: lbClient,
	}, nil
}
