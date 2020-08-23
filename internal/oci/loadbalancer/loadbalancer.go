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
	DeleteLoadBalancerByName(name string) error
}

// OCILoadBalancerController wraps logic for create,update,delete load balancers in OCI.
type OCILoadBalancerController struct {
	client        client.OCIClient
	configuration config.Config
	logger        zap.Logger
}

// NewOCILoadBalancerController will create a new OCILoadBalancerController
func NewOCILoadBalancerController(client client.OCIClient, configuration config.Config, logger zap.Logger) *OCILoadBalancerController {
	return &OCILoadBalancerController{
		client:        client,
		configuration: configuration,
		logger:        logger,
	}
}

// GetLoadBalancerByName will fetch a load balancer with a given display name if it exists
func (controller *OCILoadBalancerController) GetLoadBalancerByName(name string) (*loadbalancer.LoadBalancer, error) {
	ctx := context.Background()
	controller.logger.Sugar().Infof("Getting load balancer: %s", name)
	return controller.client.GetLoadBalancerByName(ctx, controller.configuration.Loadbalancer.Compartment, name)
}

// DeleteLoadBalancerByName will delete a load balancer with a matching display name if it exists in OCI.
func (controller *OCILoadBalancerController) DeleteLoadBalancerByName(name string) error {
	ctx := context.Background()
	controller.logger.Sugar().Infof("Deleting load balancer: %s", name)

	controller.client.DeleteLoadBalancerByName(ctx, controller.configuration.Loadbalancer.Compartment, name)
	return nil
}

// Reconcile will take an ingress and try to create or update a load balancer in OCI
func (controller *OCILoadBalancerController) Reconcile(ingress *extensions.Ingress) (*loadbalancer.LoadBalancer, error) {
	ctx := context.Background()
	definition := NewLoadBalancerDefinition(ingress, nil, []string{})

	controller.logger.Sugar().Infof("Reconciling ingress: %s", definition.Name)

	_, err := controller.client.GetLoadBalancerByName(ctx, controller.configuration.Loadbalancer.Compartment, definition.Name)
	if err != nil {
		// Perform a create action (TODO check for error being not found specifically)
		controller.createLoadBalancer(ctx, definition)
	}

	// Perform an update action
	return controller.updateLoadBalancer(ctx, definition)
}

func (controller *OCILoadBalancerController) createLoadBalancer(ctx context.Context, definition Definition) (*loadbalancer.LoadBalancer, error) {
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

	_, err := controller.client.CreateLoadBalancer(ctx, req)
	if err != nil {
		controller.logger.Sugar().Errorf("Failed to create load balancer: %s", err)
		return nil, err
	}

	return nil, nil
}

func (controller *OCILoadBalancerController) updateLoadBalancer(ctx context.Context, definition Definition) (*loadbalancer.LoadBalancer, error) {
	controller.logger.Sugar().Infof("Updating load balancer: %+s", definition.Name)
	return nil, nil
}

func (controller *OCILoadBalancerController) deleteBalancer(ctx context.Context, definition Definition) (*loadbalancer.LoadBalancer, error) {
	controller.logger.Sugar().Infof("Deleting load balancer: %+s", definition.Name)
	return nil, nil
}
