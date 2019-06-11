package main

import (
	"github.com/owainlewis/oci-kubernetes-ingress/internal/ingress/controller"

	"github.com/owainlewis/oci-kubernetes-ingress/internal/oci/client"
	conf "github.com/owainlewis/oci-kubernetes-ingress/internal/oci/config"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/settings"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

var controllerName = "oracle-cloud-infrastructure-ingress-controller"

func main() {
	// Setup logging
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Starting ingress controller")

	// Read command flags
	settings, err := settings.Load()
	if err != nil {
		logger.Sugar().Fatalf("Failed to load settings: %s", err)
	}

	// Read configuration
	c, err := conf.FromFile(settings.Config)
	if err != nil || c == nil {
		logger.Sugar().Infof("Failed to load configuration: %s", err)
	}

	// Build configuration provider
	provider, err := conf.NewConfigurationProvider(c)
	if err != nil {
	}

	// Build generic client from configuration provider
	ociClient, err := client.NewOCI(provider)
	if err != nil {
		logger.Fatal("Failed to build OCI client")
	}

	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
		Namespace: settings.Namespace,
	})
	if err != nil {
		logger.Fatal("Failed to start manager")
	}

	if err := controller.Initialize(mgr, *c, ociClient, *logger); err != nil {
		logger.Fatal("Failed to initialize controller")
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		logger.Fatal("Failed to start manager")
	}
}
