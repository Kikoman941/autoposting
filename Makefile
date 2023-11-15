REFLEX ?= github.com/cespare/reflex
GOLANGCI_LINT ?= $(GOPATH)/bin/golangci-lint
DOCKER_BIN := $(shell command -v docker 2> /dev/null)

.PHONY:

dev:
	go run $(REFLEX) -R "\\.idea|vendor|tests" -r "\\.go" -s -- sh -c "go run --race ./cmd/main.go"


lint:
	gofmt -w cmd/ internal/
	$(GOLANGCI_LINT) run --config .golangci.yml ./...

sqlc-gen:
	"$(DOCKER_BIN)" run --env-file .env --rm -v $(shell pwd):/src -w /src sqlc/sqlc -f ./internal/infrastructure/sqlc-pg/sqlc.yml generate

gen:
	go generate ./internal/presentation/graphql/