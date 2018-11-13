package ingress

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/loadbalancer"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/config"
	"k8s.io/apimachinery/pkg/util/wait"
)

var workRequestPollInterval = 5 * time.Second

// WorkRequestID is a type alias for an OCI work request OCID.
type WorkRequestID = string

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
func (svc *LoadBalancerService) CreateLoadBalancer(specification Specification) (WorkRequestID, error) {
	request := loadbalancer.CreateLoadBalancerRequest{
		CreateLoadBalancerDetails: loadbalancer.CreateLoadBalancerDetails{
			CompartmentId: common.String(specification.GetLoadBalancerCompartment()),
			DisplayName:   common.String(GetLoadBalancerUniqueName(specification.Ingress)),
			ShapeName:     common.String(specification.GetLoadBalancerShape()),
			IsPrivate:     common.Bool(specification.LoadBalancerIsPrivate()),
			SubnetIds:     specification.GetLoadBalancerSubnets(),
			// PathRouteSets:
			// BackendSets:   spec.BackendSets,
			// Listeners:     spec.Listeners,
			// Certificates:  certs,
		},
	}

	ctx := context.Background()
	response, err := svc.client.CreateLoadBalancer(ctx, request)

	return *response.OpcWorkRequestId, err
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
func (svc *LoadBalancerService) GetLoadBalancer(compartment string, name string) (*loadbalancer.LoadBalancer, error) {
	var page *string
	for {
		request := loadbalancer.ListLoadBalancersRequest{
			CompartmentId: common.String(compartment),
			DisplayName:   common.String(name),
			Page:          page,
		}

		ctx := context.Background()
		response, err := svc.client.ListLoadBalancers(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, lb := range response.Items {
			if lb.DisplayName != nil && *lb.DisplayName == name {
				return &lb, nil
			}
		}

		if page = response.OpcNextPage; page == nil {
			break
		}
	}

	return nil, fmt.Errorf("Could not find load balancer with display name %s", name)
}

// GetWorkRequest will fetch a WorkRequest from OCI for a given WorkRequest ID.
func (svc *LoadBalancerService) GetWorkRequest(ctx context.Context, id string) (*loadbalancer.WorkRequest, error) {
	resp, err := svc.client.GetWorkRequest(ctx, loadbalancer.GetWorkRequestRequest{
		WorkRequestId: common.String(id),
	})

	if err != nil {
		return nil, err
	}

	return &resp.WorkRequest, nil
}

// AwaitWorkRequest will block (poll) for a work request to reach a success or failure state.
func (svc *LoadBalancerService) AwaitWorkRequest(ctx context.Context, id string) (*loadbalancer.WorkRequest, error) {
	var wr *loadbalancer.WorkRequest
	err := wait.PollUntil(workRequestPollInterval, func() (done bool, err error) {
		glog.V(4).Infof("Polling work request: %s", id)
		workReq, err := svc.GetWorkRequest(ctx, id)
		if err != nil {
			return true, err
		}
		switch workReq.LifecycleState {
		case loadbalancer.WorkRequestLifecycleStateSucceeded:
			wr = workReq
			return true, nil
		case loadbalancer.WorkRequestLifecycleStateFailed:
			return false, fmt.Errorf("WorkRequest %q failed: %s", id, *workReq.Message)
		}
		return false, nil
	}, ctx.Done())

	return wr, err
}
