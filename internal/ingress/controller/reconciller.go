package controller

import (
	"context"

	"go.uber.org/zap"

	"k8s.io/apimachinery/pkg/api/errors"

	extensions "k8s.io/api/extensions/v1beta1"

	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/owainlewis/oci-kubernetes-ingress/internal/oci/config"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/oci/loadbalancer"
)

// Reconciler reconciles a single ingress
type Reconciler struct {
	client        client.Client
	cache         cache.Cache
	configuration config.Config
	controller    loadbalancer.Controller
	logger        zap.Logger
}

// Reconcile will reconcile the aws resources with k8s state of ingress.
func (r *Reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	r.logger.Info("Reconcile loop called")
	ctx := context.Background()
	ingress := &extensions.Ingress{}

	if err := r.cache.Get(ctx, request.NamespacedName, ingress); err != nil {
		if !errors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		r.logger.Sugar().Infof("Could not find ingress. Deleting", ingress)

		return reconcile.Result{}, nil
	}

	r.logger.Sugar().Info("Creating a new load balancer for ingress %s", ingress.Name)
	r.controller.Create(ctx, ingress)

	return reconcile.Result{}, nil

}

func (r *Reconciler) updateIngressStatus(ctx context.Context, ingress *extensions.Ingress) error {
	return r.client.Status().Update(ctx, ingress)
}
