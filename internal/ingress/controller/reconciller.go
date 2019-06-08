package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	extensions "k8s.io/api/extensions/v1beta1"

	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconciler reconciles a single ingress
type Reconciler struct {
	client client.Client
	cache  cache.Cache
}

// Reconcile will reconcile the aws resources with k8s state of ingress.
func (r *Reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	fmt.Println("Reconcile loop called")
	ctx := context.Background()
	ingress := &extensions.Ingress{}

	if err := r.cache.Get(ctx, request.NamespacedName, ingress); err != nil {
		if !errors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		if err := r.deleteIngress(ctx, request.NamespacedName); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	if err := r.reconcileIngress(ctx, request.NamespacedName, ingress); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil

}

func (r *Reconciler) reconcileIngress(ctx context.Context, ingressKey types.NamespacedName, ingress *extensions.Ingress) error {
	return nil
}

func (r *Reconciler) deleteIngress(ctx context.Context, ingressKey types.NamespacedName) error {
	return nil
}

func (r *Reconciler) updateIngressStatus(ctx context.Context, ingress *extensions.Ingress) error {
	return r.client.Status().Update(ctx, ingress)
}
