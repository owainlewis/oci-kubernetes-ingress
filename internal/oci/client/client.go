package client

import (
	"github.com/oracle/oci-go-sdk/loadbalancer"

	"github.com/oracle/oci-go-sdk/common"
)

type OCI struct {
	Loadbalancer loadbalancer.LoadBalancerClient
}

func NewOCI(provider common.ConfigurationProvider) (OCI, error) {
	lbClient, err := loadbalancer.NewLoadBalancerClientWithConfigurationProvider(provider)
	if err != nil {
		return OCI{}, err
	}

	return OCI{
		Loadbalancer: lbClient,
	}, nil
}
