.PHONY: build install test lint gomod

SHELL                = /bin/bash -o pipefail
GO_TEST_FLAGS        = -timeout 5m
GO_TEST_EXTRA_FLAGS ?=

build:
	go build ./...

test:
	go test `go list ./... | grep -v 'turbine-go\/init'` \
		$(GO_TEST_FLAGS) $(GO_TEST_EXTRA_FLAGS) \
		./...

gomod:
	go mod tidy

lint:
	golangci-lint run --timeout 5m -v

#.PHONY: generate
#generate: mockgen-install
#	go generate ./...

#mockgen-install:
#	go install github.com/golang/mock/mockgen@$(MOCKGEN_VERSION)

.PHONY: fmt
fmt: gofumpt
	gofumpt -l -w .

gofumpt:
	go install mvdan.cc/gofumpt@latest
