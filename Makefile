REFLEX ?= github.com/cespare/reflex
GOLANGCI_LINT ?= $(GOPATH)/bin/golangci-lint

.PHONY:

dev:
	go run $(REFLEX) -R "\\.idea|vendor|tests" -r "\\.go" -s -- sh -c "go run --race ./cmd/main.go"


lint:
	gofmt -w cmd/ internal/
	$(GOLANGCI_LINT) run --config .golangci.yml ./...

gen:
	go generate ./internal/presentation/graphql/