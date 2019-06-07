package main

import (
	"fmt"
	"os"

	v1beta1 "k8s.io/api/extensions/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"

	rec "github.com/owainlewis/oci-kubernetes-ingress/internal/controller"
)

func main() {

	fmt.Println("Starting...")

	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		os.Exit(1)
	}

	// Setup a new controller to reconcile ReplicaSets
	c, err := controller.New("ingress-controller", mgr, controller.Options{
		Reconciler: rec.NewOracleIngressReconciler(mgr.GetClient()),
	})
	if err != nil {
		fmt.Println("Failed to create controller")
		os.Exit(1)
	}

	// Watch ReplicaSets and enqueue ReplicaSet object key
	if err := c.Watch(&source.Kind{Type: &v1beta1.Ingress{}}, &handler.EnqueueRequestForObject{}); err != nil {
		fmt.Println("Failed to watch")
		os.Exit(1)
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		os.Exit(1)
	}
}
