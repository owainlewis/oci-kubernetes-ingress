package main

import (
	"flag"
	"log"
	"time"

	"github.com/owainlewis/oci-kubernetes-ingress/pkg/config"
	"github.com/owainlewis/oci-kubernetes-ingress/pkg/controller"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	kubeinformers "k8s.io/client-go/informers"

	"github.com/sirupsen/logrus"
)

var (
	kubeconfig = flag.String("kubeconfig", "", "Path to a kubeconfig file")
	namespace  = flag.String("namespace", "default", "Namespace to run in")
	configfile = flag.String("config", "config.yaml", "Path to the ingress controller configuration file")
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	// Load configuration
	configuration, err := config.Read(*configfile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	logrus.Debugf("Using configuration %+v", configuration)

	// Load Kubernetes client
	client, err := buildClient(*kubeconfig)
	if err != nil {
		logrus.Fatalf("Failed to create kubernetes client: %s", err)
	}

	// Init controllers
	informerFactory := kubeinformers.NewSharedInformerFactory(client, time.Second*30)
	ctrl := controller.NewOCIController(*configuration, client, *namespace, informerFactory)

	stopCh := make(chan struct{})
	go informerFactory.Start(stopCh)

	logrus.Info("Starting OCI Ingress Controller")
	ctrl.Run(1, stopCh)
}

// buildClient will construct a K8s clientset based on either local
// or in-cluster configuration depending on context
func buildClient(kubeconfig string) (kubernetes.Interface, error) {
	var config *rest.Config
	var err error
	if kubeconfig != "" {
		logrus.Debugf("Using local kubeconfig at path %s", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
