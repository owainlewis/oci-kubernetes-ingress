package main

import (
	"context"

	"github.com/owainlewis/oci-kubernetes-ingress/internal/oci/client"
	conf "github.com/owainlewis/oci-kubernetes-ingress/internal/oci/config"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/oci/loadbalancer"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/settings"
	"go.uber.org/zap"

	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
	logger.Sugar().Infof("Path %s", settings.Config)

	c, err := conf.FromFile(settings.Config)
	if err != nil || c == nil {
		logger.Sugar().Infof("Failed to load configuration: %s", err)
	}

	logger.Sugar().Infof("Config is %+v", c)

	// Build configuration provider
	provider, err := conf.NewConfigurationProvider(c)
	if err != nil {
		logger.Sugar().Fatalf("Failed to build configuration provider: %s", err)
	}

	// Build generic client from configuration provider
	ociClient, err := client.NewOCI(provider)
	if err != nil {
		logger.Fatal("Failed to build OCI client")
	}

	ingress := &extensions.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx",
		},
		Spec: extensions.IngressSpec{},
	}

	definition := loadbalancer.NewLoadBalancerDefinition(ingress)

	lbc := loadbalancer.NewOCILoadBalancerController(ociClient, *c, *logger)

	ctx := context.Background()

	lbc.Create(ctx, definition)

}
