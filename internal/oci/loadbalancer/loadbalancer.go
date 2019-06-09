package loadbalancer

import (
	"context"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"

	"github.com/owainlewis/oci-kubernetes-ingress/internal/oci/client"

	"github.com/oracle/oci-go-sdk/loadbalancer"
	extensions "k8s.io/api/extensions/v1beta1"
)

// Controller maps Kubernetes Ingress objects to OCI load balancers.
type Controller interface {
	Reconcile(ctx context.Context, ingress *extensions.Ingress) (*loadbalancer.LoadBalancer, error)
	Delete(ctx context.Context, ingressKey types.NamespacedName) error
}

type OCILoadBalancerController struct {
	client client.OCI
	logger zap.Logger
}

func NewOCILoadBalancerController(logger zap.Logger, client client.OCI) *OCILoadBalancerController {
	return &OCILoadBalancerController{
		logger: logger,
		client: client,
	}
}

func (lb *OCILoadBalancerController) Reconcile(ctx context.Context, ingress *extensions.Ingress) (*loadbalancer.LoadBalancer, error) {
	return nil, nil
}

func (lb *OCILoadBalancerController) Delete(ctx context.Context, ingressKey types.NamespacedName) error {
	return nil
}
