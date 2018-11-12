package config

import (
	"testing"
)

var validConfig = `
loadbdalancer:
  compartment: ocid1.compartment.oc1...
  subnets:
    - ""
    - ""
auwth:
  region: us-phoenix-1
  tenancy: ocid1.tenancy.oc1...
  user: ocid1.user.oc1...
  key: |
    -----BEGIN RSA PRIVATE KEY-----
    -----END RSA PRIVATE KEY-----
  fingerprint: 97:84:f7:26:a3:7b:74:...
`

func TestParseConfigInvalidConfig(t *testing.T) {
	_, err := parseConfig([]byte("Invalid config"))
	if err == nil {
		t.Fatal("expected error when given invalid config")
	}
}
func TestParseConfigValidConfig(t *testing.T) {
	_, err := parseConfig([]byte(validConfig))
	if err != nil {
		t.Fatalf("expected no error but got '%s'", err)
	}
}
