PROJECT = oci-ingress
REGISTRY ?= iad.ocir.io/oracle
IMAGE := $(REGISTRY)/$(PROJECT)

.PHONY: test
test:
	go test ./...

.PHONY: run
run:
	go run cmd/main.go -kubeconfig=/Users/owainlewis/.kube/config -config=config.yml

.PHONY: deps
deps:
	dep ensure -vendor-only
