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

// OCILoadBalancerController wraps logic for create,update,delete load balancers in OCI.
type OCILoadBalancerController struct {
	client        client.OCI
	configuration config.Config
	logger        zap.Logger
}

// NewOCILoadBalancerController will create a new OCILoadBalancerController
func NewOCILoadBalancerController(client client.OCI, configuration config.Config, logger zap.Logger) *OCILoadBalancerController {
	return &OCILoadBalancerController{
		client:        client,
		configuration: configuration,
		logger:        logger,
	}
}

// Reconcile will take an ingress and try to create or update a load balancer in OCI
func (controller *OCILoadBalancerController) Reconcile(ingress *extensions.Ingress) (*loadbalancer.LoadBalancer, error) {
	controller.logger.Sugar().Infof("Reconciling ingress: %s/%s", ingress.Namespace, ingress.Name)

	ctx := context.Background()
	definition := NewLoadBalancerDefinition(ingress)
	return controller.createLoadBalancer(ctx, definition)

}

func (controller *OCILoadBalancerController) createLoadBalancer(ctx context.Context, definition LoadBalancerDefinition) (*loadbalancer.LoadBalancer, error) {
	controller.logger.Sugar().Infof("Creating load balancer from definition: %+v", definition)
	details := loadbalancer.CreateLoadBalancerDetails{
		CompartmentId: common.String(controller.configuration.Loadbalancer.Compartment),
		DisplayName:   common.String(definition.Name),
		ShapeName:     common.String(definition.Shape),
		IsPrivate:     common.Bool(false),
		SubnetIds:     controller.configuration.Loadbalancer.Subnets,
	}

	req := loadbalancer.CreateLoadBalancerRequest{
		CreateLoadBalancerDetails: details,
	}

	_, err := controller.client.Loadbalancer.CreateLoadBalancer(ctx, req)
	if err != nil {
		controller.logger.Sugar().Errorf("Failed to create load balancer: %s", err)
		return nil, err
	}

	return nil, nil
}
