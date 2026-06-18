VERSION  := $(shell git describe --tags --always 2>/dev/null || echo dev)
LDFLAGS  := -s -w -X main.Version=$(VERSION)
BINARY   := grimoire
BIN_DIR  := bin
DIST_DIR := dist

PLATFORMS := \
	linux/amd64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64

.PHONY: build dev install test vet fmt tidy lint clean dist

build: ## Build local binary to bin/
	@mkdir -p $(BIN_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) .

dev: ## Run without building (go run .) — pass args with ARGS="..."
	go run . $(ARGS)

install: ## Install to GOPATH/bin
	go install -ldflags="$(LDFLAGS)" .

test: ## Run all tests
	go test ./...

vet: ## Run go vet
	go vet ./...

fmt: ## Format source files
	gofmt -l -w .

tidy: ## Sync go.mod and go.sum
	go mod tidy

lint: ## Run golangci-lint
	golangci-lint run ./...

clean: ## Remove build artifacts
	rm -rf $(BIN_DIR) $(DIST_DIR)

dist: ## Cross-compile all platforms to dist/
	@mkdir -p $(DIST_DIR)
	$(foreach platform,$(PLATFORMS), \
		$(eval GOOS   := $(word 1,$(subst /, ,$(platform)))) \
		$(eval GOARCH := $(word 2,$(subst /, ,$(platform)))) \
		$(eval EXT    := $(if $(filter windows,$(GOOS)),.exe,)) \
		$(eval OUT    := $(DIST_DIR)/$(BINARY)-$(GOOS)-$(GOARCH)$(EXT)) \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="$(LDFLAGS)" -o $(OUT) . && \
		echo "  built: $(OUT)" ; \
	)

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
