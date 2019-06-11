package settings

import (
	"flag"

	apiv1 "k8s.io/api/core/v1"
)

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

func Load() (*Settings, error) {
	settings := &Settings{
		Namespace: defaultNamespace,
		Config:    defaultConfigPath,
	}

	settings.bindAll()

	flag.Parse()

	return settings, nil
}
