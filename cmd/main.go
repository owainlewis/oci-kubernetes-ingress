package main

import (
	"os"

	"github.com/owainlewis/oci-kubernetes-ingress/internal/ingress/controller"

	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

var controllerName = "oracle-cloud-infrastructure-ingress-controller"

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Starting ingress controller")

	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
		//		Namespace: "default",
	})
	if err != nil {
		os.Exit(1)
	}

	if err := controller.Initialize(mgr); err != nil {
		logger.Fatal("Failed to initialize controller")
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		logger.Fatal("Failed to start manager")
	}
}
