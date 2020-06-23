package controller

import (
	"context"

	"go.uber.org/zap"

	"k8s.io/apimachinery/pkg/api/errors"

	corev1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"

	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/owainlewis/oci-kubernetes-ingress/internal/common"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/ingress/controller/store"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/oci/config"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/oci/loadbalancer"

	oci "github.com/oracle/oci-go-sdk/loadbalancer"
)

// Reconciler reconciles a single ingress
type Reconciler struct {
	client        client.Client
	cache         cache.Cache
	store         store.Store
	configuration config.Config
	controller    loadbalancer.Controller
	logger        zap.Logger
}

func (r *Reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	r.logger.Sugar().Infof("Reconcile loop called for %s", request.NamespacedName)

	ctx := context.Background()
	ingress := &extensions.Ingress{}
	if err := r.cache.Get(ctx, request.NamespacedName, ingress); err != nil {
		if !errors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		r.logger.Sugar().Infof("Could not find ingress to reconcile: %s", request.NamespacedName)
		loadBalancerName := common.GetLoadBalancerName(request.Namespace, request.Name)
		r.deleteIngressByName(ctx, loadBalancerName)

		return reconcile.Result{}, nil
	}

	r.logger.Sugar().Infof("Reconciling ingress: %s", request.NamespacedName)
	r.reconcileIngress(ctx, ingress)

	return reconcile.Result{}, nil
}

func (r *Reconciler) deleteIngressByName(ctx context.Context, name string) {
	r.logger.Sugar().Infof("Deleting ingress: %s", name)
	r.controller.DeleteLoadBalancerByName(name)
}

func (r *Reconciler) reconcileIngress(ctx context.Context, ingress *extensions.Ingress) error {
	lb, err := r.controller.Reconcile(ingress)
	if err != nil {
		return err
	}
	if err := r.updateIngressStatus(ctx, ingress, lb); err != nil {
		return err
	}

	return nil
}

func (r *Reconciler) updateIngressStatus(ctx context.Context, ingress *extensions.Ingress, lb *oci.LoadBalancer) error {
	ingress.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{
		{
			Hostname: "todo",
		},
	}

	return r.client.Status().Update(ctx, ingress)
}
