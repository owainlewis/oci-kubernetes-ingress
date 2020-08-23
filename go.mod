module github.com/owainlewis/oci-kubernetes-ingress

go 1.15

require (
	github.com/oracle/oci-go-sdk v24.1.0+incompatible
	go.uber.org/zap v1.15.0
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	sigs.k8s.io/controller-runtime v0.6.2
)
