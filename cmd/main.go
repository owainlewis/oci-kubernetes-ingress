package main

import (
	"flag"
	"time"

	"github.com/golang/glog"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/config"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/context"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/controller"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "",
		"Path to a kubeconfig file")

	namespace = flag.String("namespace", "default",
		"Namespace to run in")

	configfile = flag.String("config", "cloud-provider.yaml",
		"Path to the OCI ingress controller configuration file.")

	interval = flag.Duration("interval", 30*time.Minute,
		"The reconcile interval for Kubernetes informers")
)

func main() {
	flag.Parse()

	// Load OCI configuration file
	configuration, err := loadAndValidateConfiguration(*configfile)
	if err != nil {
		glog.Fatalf("Invalid or absent configuration: %v", err)
	}

	// Load Kubernetes client
	kubeClient, err := buildK8sClient(*kubeconfig)
	if err != nil {
		glog.Fatalf("Failed to create kubernetes client: %v", err)
	}

	context := context.NewControllerContext(kubeClient, *namespace, *interval)
	stopCh := make(chan struct{})
	ctrl := controller.NewOCIController(*configuration, context, stopCh)

	ctrl.Run()
}

// buildK8sClient will construct a K8s client based on either local
// or in-cluster configuration.
func buildK8sClient(kubeconfig string) (kubernetes.Interface, error) {
	var config *rest.Config
	var err error
	if kubeconfig != "" {
		glog.V(4).Infof("Using local kubeconfig at path %s", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

// loadAndValidateConfiguration will read and validate a configuration file
// from disk.
func loadAndValidateConfiguration(filepath string) (*config.Config, error) {
	configuration, err := config.FromFile(filepath)
	if err != nil {
		return nil, err
	}

	err = configuration.Validate()
	if err != nil {
		return nil, err
	}

	return configuration, nil
}
