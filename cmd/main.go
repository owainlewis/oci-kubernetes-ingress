package main

import (
	"flag"

	"github.com/owainlewis/oci-kubernetes-ingress/internal/ingress/controller"
	apiv1 "k8s.io/api/core/v1"

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

	settings, err := loadSettings()
	if err != nil {
		logger.Fatal("Failed to load settings")
	}

	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
		Namespace: settings.Namespace,
	})
	if err != nil {
		logger.Fatal("Failed to start manager")
	}

	if err := controller.Initialize(mgr); err != nil {
		logger.Fatal("Failed to initialize controller")
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		logger.Fatal("Failed to start manager")
	}
}

const (
	defaultNamespace = apiv1.NamespaceAll
)

// Settings defines common settings for the ingress controller
type Settings struct {
	Namespace string
}

func (settings *Settings) bindFlags(fs *flag.FlagSet) {
	fs.StringVar(&settings.Namespace, "namespace", defaultNamespace,
		`Namespace sets the controller watch namespace for updates to Kubernetes objects.
		Defaults to all namespaces if not set.`)
}

func loadSettings() (*Settings, error) {
	settings := &Settings{
		Namespace: defaultNamespace,
	}

	fs := flag.NewFlagSet("", flag.ExitOnError)
	settings.bindFlags(fs)

	return settings, nil
}
