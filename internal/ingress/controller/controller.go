package controller

import (
	"fmt"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"

	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/owainlewis/oci-kubernetes-ingress/internal/ingress/controller/handlers"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/oci/client"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/oci/config"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/oci/loadbalancer"
)

// Initialize ...
func Initialize(mgr manager.Manager, configuration config.Config, client client.OCIClient, logger zap.Logger) error {
	loadbalancerController := loadbalancer.NewOCILoadBalancerController(client, configuration, logger)

	reconciler, err := newReconciler(mgr, configuration, loadbalancerController, logger)
	if err != nil {
		return err
	}

	c, err := controller.New("oci-ingress-controller", mgr, controller.Options{Reconciler: reconciler})
	if err != nil {
		return err
	}

	if err := watchClusterEvents(c, mgr.GetCache()); err != nil {
		return fmt.Errorf("failed to watch cluster events: %v", err)
	}

	return nil
}

func newReconciler(mgr manager.Manager, configuration config.Config, lbc loadbalancer.Controller, logger zap.Logger) (reconcile.Reconciler, error) {

	return &Reconciler{
		client:        mgr.GetClient(),
		cache:         mgr.GetCache(),
		configuration: configuration,
		controller:    lbc,
		logger:        logger,
	}, nil
}

func watchClusterEvents(c controller.Controller, cache cache.Cache) error {
	// Watch Ingress objects for changes (Create, Update, Delete)
	if err := c.Watch(&source.Kind{Type: &extensions.Ingress{}}, &handlers.EnqueueRequestsForIngressEvent{}); err != nil {
		return err
	}

	// Watch Node objects for changes (Create, Delete)
	if err := c.Watch(&source.Kind{Type: &corev1.Node{}}, &handlers.EnqueueRequestsForNodeEvent{
		Cache: cache,
	}); err != nil {
		return err
	}

	return nil
}
