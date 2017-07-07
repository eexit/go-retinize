# Env Variables
# =============================================================================================
ROOT_DIR?=$(PWD)
BINARY_PATH=$(ROOT_DIR)/bin/retinize

GO_VERSION=1.8
SHELL=/bin/bash

# Rules
# =============================================================================================
.PHONY: build test
.DEFAULT: usage

usage:
	@echo '+---------------------------------------------------------------------------------------+'
	@echo '| Retinize Make Usage                                                                   |'
	@echo '+---------------------------------------------------------------------------------------+'
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "|- \033[33m%-15s\033[0m -> %s\n", $$1, $$2}'

check-go: ## Checks Go configuration
	@go version >/dev/null 2>&1 || { echo >&2 "Go is required.  Aborting."; exit 1; }
	@go version | grep -q "^go version go$(GO_VERSION)" 2>&1 || { printf >&2 "\033[31mGo $(GO_VERSION) is required.\033[0m\nAborting.\n"; exit 1; }
	@printf "\033[32mGOPATH\033[0m   = $(GOPATH)\n\033[32mROOT_DIR\033[0m = $(ROOT_DIR)\n"

install: ## Installs dependencies
	@glide --version >/dev/null 2>&1 || { echo >&2 "Glide is required.  Aborting."; exit 1; }
	@glide install
	@golint --help >/dev/null 2>&1 || echo "Installing golint…" && go get -u github.com/golang/lint/golint

update: ## Updates dependencies (using glide)
	@glide --version >/dev/null 2>&1 || { echo >&2 "Glide is required.  Aborting."; exit 1; }
	@glide update

build: check-go ## Compiles and pack the Go binary for use via the docker container
	@echo "Compiling..." && go build -o $(BINARY_PATH) $(ROOT_DIR)/main.go

test: ## Launchs unit test suite
	@bin/coverage --html

lint: ## Runs golint on project files
	@go list github.com/eexit/go-retinize/... | grep -v /vendor | xargs golint -set_exit_status && echo "All good!"

vet: ## Runs go tool vet on project files
	@find . -type f -name "*.go" -not -path "./vendor/*" | xargs -I {} go tool vet {} && echo "All good!"

