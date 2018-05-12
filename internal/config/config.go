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

// Config defines the configuration needed for the OCI ingress controller
type Config struct {
	Auth AuthConfig `yaml:"auth"`
}

// Read will try and parse a config object a file
func Read(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
