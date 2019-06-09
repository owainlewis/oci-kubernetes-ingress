package client

import (
	"github.com/oracle/oci-go-sdk/loadbalancer"

	"github.com/oracle/oci-go-sdk/common"
)

type OCI struct {
	loadbalancer loadbalancer.LoadBalancerClient
}

func NewOCI(provider common.ConfigurationProvider) (OCI, error) {
	return OCI{}, nil
}
