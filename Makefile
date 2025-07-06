##@ Top-Level

.PHONY: all
all: deps fmt build run test ## Run the full generation pipeline (deps, fmt, build, run)

##@ Basic Targets

.PHONY: run
run: ## Run the alias generator
	go run .

.PHONY: build
build: ## Build the alias generator binary
	go build -o kubectl-aliases

.PHONY: clean
clean: ## Remove generated files
	rm -f kubectl-aliases kubectl_aliases

.PHONY: fmt
fmt: ## Format all Go files
	go fmt ./...

.PHONY: start-kind
start-kind: ## Start a kind cluster if not already running
	@if ! kubectl cluster-info --context kind-kind >/dev/null 2>&1; then \
		echo "No kind cluster found. Starting a new kind cluster..."; \
		kind create cluster; \
	else \
		echo "Kind cluster is already running."; \
	fi

.PHONY: test
test: start-kind ## Run tests (after ensuring kind is running)
	go test ./...

##@ Dependency Management

.PHONY: deps
deps: ## Ensure dependencies are tidy and verified
	go mod tidy
	go mod verify

.PHONY: upgrade
upgrade: ## Upgrade all Go module dependencies
	go get -u ./...
	go mod tidy

##@ Help

.PHONY: help
help: ## Display Makefile help information for all actions
	@awk 'BEGIN {FS = ":.*##"; \
		printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} \
		/^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' \
		$(MAKEFILE_LIST)
