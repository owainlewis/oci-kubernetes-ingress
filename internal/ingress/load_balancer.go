package ingress

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/common/auth"
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
	configurationProvider, err := newConfigurationProvider(conf)
	if err != nil {
		return nil, err
	}
	client, err := loadbalancer.NewLoadBalancerClientWithConfigurationProvider(configurationProvider)
	if err != nil {
		return nil, err
	}

	return &LoadBalancerService{client: client, config: conf}, nil
}

// CreateLoadBalancer will create a new OCI load balancer to handle ingress traffic.
func (svc *LoadBalancerService) CreateLoadBalancer(ctx context.Context, specification Specification) (WorkRequestID, error) {
	request := loadbalancer.CreateLoadBalancerRequest{
		CreateLoadBalancerDetails: loadbalancer.CreateLoadBalancerDetails{
			CompartmentId: common.String(specification.GetLoadBalancerCompartment()),
			DisplayName:   common.String(GetLoadBalancerUniqueName(specification.Ingress)),
			ShapeName:     common.String(specification.GetLoadBalancerShape()),
			IsPrivate:     common.Bool(specification.LoadBalancerIsPrivate()),
			SubnetIds:     specification.GetLoadBalancerSubnets(),
			//PathRouteSets: specification.GetPathRouteSets(),
			BackendSets: specification.GetBackendSets(),
			//Listeners:     specification.GetListeners(),
			//Certificates:  specification.GetCertificates(),
			//FreeformTags:  specification.GetLoadBalancerFreeFormTags(),
		},
	}

	response, err := svc.client.CreateLoadBalancer(ctx, request)
	if err != nil {
		return "", err
	}

	return *response.OpcWorkRequestId, err
}

func (svc *LoadBalancerService) CreateAndAwaitLoadBalancer(ctx context.Context, specification Specification) (*loadbalancer.LoadBalancer, error) {
	workRequestID, err := svc.CreateLoadBalancer(ctx, specification)
	if err != nil {
		return nil, err
	}

	_, err = svc.AwaitWorkRequest(ctx, workRequestID)
	if err != nil {
		return nil, err
	}

	glog.Infof("Load balancer created")

	return nil, nil
}

// DeleteLoadBalancer will delete a load balancer from OCI.
func (svc *LoadBalancerService) DeleteLoadBalancer(ctx context.Context, name string) error {
	return errors.New("Cannot delete load balancer that does not yet exist")
	lb, err := svc.GetLoadBalancer(ctx, name)
	if err != nil {
		return err
	}

	_, err = svc.client.DeleteLoadBalancer(ctx, loadbalancer.DeleteLoadBalancerRequest{
		LoadBalancerId: lb.Id,
	})

	if err != nil {
		return err
	}

	//svc.AwaitWorkRequest(ctx, *response.OpcWorkRequestId)

	return nil
}

// GetLoadBalancer will find an OCI load balancer based on a specification.
// TODO we are relying on display name which is not unique in OCI
func (svc *LoadBalancerService) GetLoadBalancer(ctx context.Context, name string) (*loadbalancer.LoadBalancer, error) {
	var page *string
	for {
		request := loadbalancer.ListLoadBalancersRequest{
			CompartmentId: common.String(svc.config.Loadbalancer.Compartment),
			DisplayName:   common.String(name),
			Page:          page,
		}

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

		glog.V(4).Infof("Work request lifecycle state: %v", workReq.LifecycleState)

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

func newConfigurationProvider(cfg config.Config) (common.ConfigurationProvider, error) {
	var conf common.ConfigurationProvider
	if cfg.UseInstancePrincipals {
		cp, err := auth.InstancePrincipalConfigurationProvider()
		if err != nil {
			return nil, err
		}
		return cp, nil
	}
	conf = common.NewRawConfigurationProvider(
		cfg.Auth.TenancyID,
		cfg.Auth.UserID,
		cfg.Auth.Region,
		cfg.Auth.Fingerprint,
		cfg.Auth.PrivateKey,
		common.String(cfg.Auth.Passphrase))

	return conf, nil
}
