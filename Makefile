export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1
GOOS=$(shell go env GOOS)
VERSION=$(shell git describe --tags --always)
REVISION=$(shell git rev-parse HEAD)$(shell if ! git diff --no-ext-diff --quiet --exit-code; then echo .m; fi)

dep:
	go mod tidy
	go get golang.org/x/tools/cmd/goimports
	go install golang.org/x/tools/cmd/goimports

generate:
	rm -rf ./internal/mocks/*
	PATH=$(GOPATH)/bin:${PATH}
	go install github.com/golang/mock/mockgen/...@v1.6.0
	go generate -x ./internal/...

build:
	CGO_ENABLED=0 GOOS=${GOOS} go build -ldflags "-X=main.Revision=${REVISION} -X=main.Version=${VERSION}" -o ./sqsdumper ./cmd/main.go

lint:
	revive --config=revive.toml --formatter=unix ./...

fmt:
	go fmt ./...
	goimports -local andboson   -w .

test:
	go test -tags="unit" -coverprofile=coverage.out -covermode=count ./...
	go tool cover -func=coverage.out

all: generate dep test build

