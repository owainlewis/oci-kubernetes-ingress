.PHONY: run
run:
	go run cmd/main.go -kubeconfig=/Users/owainlewis/.kube/config

deps:
	dep ensure -vendor-only
