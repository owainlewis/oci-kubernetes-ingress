package controller

import (
	"fmt"

	"github.com/owainlewis/oci-kubernetes-ingress/internal/ingress/controller/handlers"
	extensions "k8s.io/api/extensions/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func Initialize(mgr manager.Manager) error {
	reconciler, err := newReconciler(mgr)
	if err != nil {
		return err
	}

	c, err := controller.New("oracle-cloud-ingress-controller", mgr, controller.Options{Reconciler: reconciler})
	if err != nil {
		return err
	}

	if err := watchClusterEvents(c, ""); err != nil {
		return fmt.Errorf("failed to watch cluster events due to %v", err)
	}

	return nil
}

func newReconciler(mgr manager.Manager) (reconcile.Reconciler, error) {
	return &Reconciler{
		client: mgr.GetClient(),
		cache:  mgr.GetCache(),
	}, nil
}

func watchClusterEvents(c controller.Controller, ingressClass string) error {
	if err := c.Watch(&source.Kind{Type: &extensions.Ingress{}}, &handlers.EnqueueRequestsForIngressEvent{
		IngressClass: ingressClass,
	}); err != nil {
		return err
	}

	return nil
}
