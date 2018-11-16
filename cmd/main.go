package main

import (
	"flag"
	"log"
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
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	// Load configuration
	configuration, err := config.FromFile(*configfile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	err = configuration.Validate()
	if err != nil {
		glog.Fatalf("Invalid configuration: %s", err)
	}

	// Load Kubernetes client
	kubeClient, err := buildClient(*kubeconfig)
	if err != nil {
		glog.Fatalf("Failed to create kubernetes client: %s", err)
	}

	context := context.NewControllerContext(kubeClient, *namespace, 30*time.Second)
	stopCh := make(chan struct{})
	ctrl := controller.NewOCIController(*configuration, context, stopCh)

	ctrl.Run()
}

// buildClient will construct a K8s clientset based on either local
// or in-cluster configuration depending on context
func buildClient(kubeconfig string) (kubernetes.Interface, error) {
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
