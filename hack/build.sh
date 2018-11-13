#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

go build \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o bin/oci-kubernetes-ingress cmd/main.go
