package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

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
	Compartment string `yaml:"compartment"`
}

// Config defines the configuration needed for the OCI ingress controller
type Config struct {
	//Loadbalancer LoadbalancerConfig `yaml:"loadbalancer"`
	Auth AuthConfig `yaml:"auth"`
}

// ParseConfig will parse the contents of a file into a Config object
func parseConfig(c []byte) (*Config, error) {
	cfg := &Config{}
	if err := yaml.Unmarshal(c, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Read will try and parse a config object a file
func Read(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return parseConfig(b)
}
