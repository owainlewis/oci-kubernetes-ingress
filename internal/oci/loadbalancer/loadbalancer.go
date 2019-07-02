package loadbalancer

import (
	"context"

	"github.com/owainlewis/oci-kubernetes-ingress/internal/oci/client"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/oci/config"
	"go.uber.org/zap"
	extensions "k8s.io/api/extensions/v1beta1"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/loadbalancer"
)

// Controller maps Kubernetes Ingress objects to OCI load balancers.
type Controller interface {
	Reconcile(ingress *extensions.Ingress) (*loadbalancer.LoadBalancer, error)
}

type OCILoadBalancerController struct {
	client        client.OCI
	configuration config.Config
	logger        zap.Logger
}

func NewOCILoadBalancerController(client client.OCI, configuration config.Config, logger zap.Logger) *OCILoadBalancerController {
	return &OCILoadBalancerController{
		client:        client,
		configuration: configuration,
		logger:        logger,
	}
}

func (controller *OCILoadBalancerController) Reconcile(ingress *extensions.Ingress) (*loadbalancer.LoadBalancer, error) {
	controller.logger.Sugar().Infof("Reconciling ingress: %s/%s", ingress.Namespace, ingress.Name)
	return nil, nil
}

func (lb *OCILoadBalancerController) createLoadBalancer(ctx context.Context, definition LoadBalancerDefinition) (*loadbalancer.LoadBalancer, error) {
	lb.logger.Sugar().Infof("Creating load balancer from definition: %+v", definition)

	details := loadbalancer.CreateLoadBalancerDetails{
		CompartmentId: common.String(lb.configuration.Loadbalancer.Compartment),
		DisplayName:   common.String(definition.Name),
		ShapeName:     common.String(definition.Shape),
		IsPrivate:     common.Bool(false),
		SubnetIds:     lb.configuration.Loadbalancer.Subnets,
	}

	req := loadbalancer.CreateLoadBalancerRequest{
		CreateLoadBalancerDetails: details,
	}

	_, err := lb.client.Loadbalancer.CreateLoadBalancer(ctx, req)
	if err != nil {
		lb.logger.Sugar().Errorf("Failed to provision load balancer: %s", err)
		return nil, err
	}

	return nil, nil
}
