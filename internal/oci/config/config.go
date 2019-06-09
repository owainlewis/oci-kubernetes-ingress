package config

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/common/auth"
	"gopkg.in/yaml.v2"
)

// Configuration settings for OCI authentication, load balancers etc

// AuthConfig holds the configuration required for communicating with the OCI
// API.
type AuthConfig struct {
	Region      string `yaml:"region"`
	TenancyID   string `yaml:"tenancy"`
	UserID      string `yaml:"user"`
	PrivateKey  string `yaml:"key"`
	Fingerprint string `yaml:"fingerprint"`
	Passphrase  string `yaml:"passphrase"`
}

// LoadbalancerConfig ...
type LoadbalancerConfig struct {
	Compartment string   `yaml:"compartment"`
	Subnets     []string `yaml:"subnets"`
}

// Config defines the configuration needed for the OCI ingress controller.
type Config struct {
	Auth                  AuthConfig         `yaml:"auth"`
	UseInstancePrincipals bool               `yaml:"useInstancePrincipals"`
	Loadbalancer          LoadbalancerConfig `yaml:"loadbalancer"`
}

// Validate performs basic structural validation on the configuration struct.
func (c Config) Validate() error {
	errs := []string{}

	if c.Loadbalancer.Compartment == "" {
		errs = append(errs, "Compartment must be declared for loadbalancer")
	}

	if len(c.Loadbalancer.Subnets) < 2 {
		errs = append(errs, "At least two subnets are required")
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(strings.Join(errs, ", "))
}

// Parse will parse a bytestring representation of a configuration struct into a Config object.
func Parse(c []byte) (*Config, error) {
	cfg := &Config{}
	if err := yaml.Unmarshal(c, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// FromFile will try and parse a config object into a file.
func FromFile(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return Parse(b)
}

// NewConfigurationProvider takes a cloud provider config file and returns an OCI ConfigurationProvider
// to be consumed by the OCI SDK.
func NewConfigurationProvider(cfg *Config) (common.ConfigurationProvider, error) {
	var conf common.ConfigurationProvider
	if cfg != nil {
		err := cfg.Validate()
		if err != nil {
			return nil, err
		}

		if cfg.UseInstancePrincipals {
			cp, err := auth.InstancePrincipalConfigurationProvider()
			if err != nil {
				return nil, err
			}
			return cp, nil
		}

		conf = common.NewRawConfigurationProvider(
			cfg.Auth.TenancyID,
			cfg.Auth.UserID,
			cfg.Auth.Region,
			cfg.Auth.Fingerprint,
			cfg.Auth.PrivateKey,
			common.String(cfg.Auth.Passphrase))

	} else {
		conf = common.DefaultConfigProvider()
	}

	return conf, nil
}
