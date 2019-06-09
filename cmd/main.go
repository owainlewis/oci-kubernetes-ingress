package main

import (
	"flag"

	"github.com/owainlewis/oci-kubernetes-ingress/internal/ingress/controller"
	apiv1 "k8s.io/api/core/v1"

	"github.com/owainlewis/oci-kubernetes-ingress/internal/oci/client"
	ociconfig "github.com/owainlewis/oci-kubernetes-ingress/internal/oci/config"
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
	settings, err := loadSettings()
	if err != nil {
		logger.Sugar().Fatalf("Failed to load settings: %s", err)
	}

	// Set OCI client
	cfg, err := ociconfig.FromFile(settings.Config)
	if err != nil {
		logger.Sugar().Infof("Failed to load configuration: %s", err)
	}

	logger.Sugar().Infof("Configuration: %+v", cfg)

	// Build configuration provider
	provider, err := ociconfig.NewConfigurationProvider(cfg)
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

	if err := controller.Initialize(mgr, ociClient, logger); err != nil {
		logger.Fatal("Failed to initialize controller")
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		logger.Fatal("Failed to start manager")
	}
}

const (
	defaultNamespace  = apiv1.NamespaceAll
	defaultConfigPath = "/etc/oci/config.yaml"
)

// Settings defines common settings for the ingress controller
type Settings struct {
	Namespace string
	Config    string
}

func (settings *Settings) bindAll() {
	flag.StringVar(&settings.Namespace, "namespace", defaultNamespace,
		`Namespace sets the controller watch namespace for updates to Kubernetes objects.
		Defaults to all namespaces if not set.`)
	flag.StringVar(&settings.Config, "config", defaultConfigPath,
		`The path to an OCI config yaml file containing auth credentials and configuration.
		Defaults to /etc/oci/config.yaml`)
}

func loadSettings() (*Settings, error) {
	settings := &Settings{
		Namespace: defaultNamespace,
		Config:    defaultConfigPath,
	}

	settings.bindAll()

	flag.Parse()

	return settings, nil
}
