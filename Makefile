SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

ifneq (3.82,$(firstword $(sort $(MAKE_VERSION) 3.82)))
  $(error This Make does not support .ONESHELL, use GNU Make 3.82 and newer)
endif

.DEFAULT_GOAL :=help

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= 0
BIN_SUFFIX := $(if $(filter windows,$(GOOS)),.exe,)
BIN_ARCH := $(GOOS)-$(GOARCH)

GITHUB_REF ?= dev
GIT_REF := $(shell echo "$(notdir $(GITHUB_REF))" | tr -cd '[:alnum:]._-')
GITHUB_SHA ?= dev
GIT_SHA ?= $(GITHUB_SHA)

BIN_DIR ?= bin
CMD_DIR ?= cmd
RELEASE_DIR ?= release-artifacts

BIN := $(BIN_DIR)/kubent-$(BIN_ARCH)$(BIN_SUFFIX)
ARTIFACT := $(RELEASE_DIR)/kubent-$(GIT_REF)-$(BIN_ARCH).tar.gz

SRC := $(wildcard $(CMD_DIR)/kubent/*.go) $(wildcard pkg/*/*.go)

ifeq (, $(shell command -v go))
  $(error "No go binary in PATH")
endif

## Clean, build and pack
all: clean build release-artifacts
.PHONY: all

## Prints list of tasks
help:
	@awk 'BEGIN {FS=":"} /^## .*/,/^[a-zA-Z0-9_-]+:/ { if ($$0 ~ /^## /) { desc=substr($$0, 4) } else { printf "\033[36m%-30s\033[0m %s\n", $$1, desc } }' Makefile | sort
.PHONY: help

## Build binary
build: $(BIN)
.PHONY: build

$(BIN): $(SRC) go.mod go.sum
	mkdir -p $(BIN_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build -ldflags="-s -w -X main.version=$(GIT_REF) -X main.gitSha=$(GIT_SHA)" -o "$@" ./$(CMD_DIR)/kubent

## Create release artifacts
release-artifacts: $(ARTIFACT)
.PHONY: release-artifacts

$(ARTIFACT): $(BIN)
	mkdir -p $(RELEASE_DIR)
	staging=$$(mktemp -d)
	cp $(BIN) "$$staging/kubent$(BIN_SUFFIX)"
	tar -czf "$@" -C "$$staging" "kubent$(BIN_SUFFIX)"
	rm -rf "$$staging"

## Run Go tests
test: test-fmt
	go test -v -coverprofile fmtcoverage.html ./...
.PHONY: test

## Run go and opa fmt checks
test-fmt:
	test -z "$$(go run github.com/open-policy-agent/opa@v0.70.0 fmt -l pkg/rules/rego/*)"
	test -z "$$(go fmt ./...)"
.PHONY: test-fmt

## Clean build artifacts
clean:
	rm -rf $(BIN_DIR) $(RELEASE_DIR)
.PHONY: clean
