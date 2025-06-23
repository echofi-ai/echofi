#!/usr/bin/make -f

DOCKER := $(shell which docker)

BUILD_FLAGS := -ldflags '$(ldflags)' -gcflags="all=-N -l"

install: go.sum
		go install $(BUILD_FLAGS) ./cmd/echofid

build:
	go build -o build/echofid ./cmd/echofid

###############################################################################
###                                Test                                 	###
###############################################################################
TEST_FLAGS := -v

test:
	go test $(TEST_FLAGS) ./...


###############################################################################
###                                Linting                                  ###
###############################################################################
golangci_lint_cmd=golangci-lint
golangci_version=v1.60.1

lint:
	@echo "--> Running linter"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@$(golangci_lint_cmd) run --timeout=10m
