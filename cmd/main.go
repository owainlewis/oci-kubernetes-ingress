package main

import (
	"flag"
	"time"

	"github.com/golang/glog"
	"github.com/owainlewis/oci-kubernetes-ingress/pkg/config"
	"github.com/owainlewis/oci-kubernetes-ingress/pkg/controller"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	kubeinformers "k8s.io/client-go/informers"

	log "github.com/sirupsen/logrus"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "Path to a kubeconfig file")
	namespace := flag.String("namespace", "default", "Namespace to run in")
	configfile := flag.String("config", "/opt/config/oci/config.yml", "Path to the ingress controller configuration file")

	flag.Set("logtostderr", "true")
	flag.Parse()

	_, err := config.Read(*configfile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	client, err := buildClient(*kubeconfig)
	if err != nil {
		glog.Fatalf("Failed to create kubernetes client: %s", err)
	}

	informerFactory := kubeinformers.NewSharedInformerFactory(client, time.Second*30)

	ctrl := controller.NewOCIController(client, *namespace, informerFactory)

	stopCh := make(chan struct{})

	go informerFactory.Start(stopCh)

	glog.Info("Starting Controller")
	ctrl.Run(1, stopCh)
}

func buildClient(kubeconfig string) (kubernetes.Interface, error) {
	var config *rest.Config
	var err error
	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
