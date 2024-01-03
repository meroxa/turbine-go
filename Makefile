.PHONY: build install test lint gomod

SHELL                = /bin/bash -o pipefail
GO_TEST_FLAGS        = -timeout 5m
GO_TEST_EXTRA_FLAGS ?=

.PHONY: build 
build:
	go build ./...

.PHONY: test 
test:
	go test `go list ./... | grep -v 'turbine-go\/init'` \
		$(GO_TEST_FLAGS) $(GO_TEST_EXTRA_FLAGS) \
		./...

.PHONY: lint
lint:
	golangci-lint run --timeout 5m -v

.PHONY: fmt
fmt: gofumpt
	gofumpt -l -w .

.PHONY: gofumpt
gofumpt:
	go install mvdan.cc/gofumpt@latest
