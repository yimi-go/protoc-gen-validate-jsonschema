# Build all by default, even if it's not first
.DEFAULT_GOAL := help

# ==============================================================================
# Build options

ROOT_PACKAGE=github.com/yimi-go/protoc-gen-validate-jsonschema
VERSION_PACKAGE=github.com/yimi-go/protoc-gen-validate-jsonschema

# ==============================================================================
# Includes

include build/make-rules/common.mk # make sure include common.mk at the first include line
include build/make-rules/golang.mk
include build/make-rules/release.mk
include build/make-rules/tools.mk

# ==============================================================================
# Usage

define USAGE_OPTIONS

Options:
  DEBUG            Whether to generate debug symbols. Default is 0.
  VERSION          The version information compiled into binaries.
                   The default is obtained from gsemver or git.
  COVERAGE         Minimum test coverage. Default is 60.
  V                Set to 1 enable verbose build. Default is undefined.
endef
export USAGE_OPTIONS

# ==============================================================================
# Targets

## git.hooks: Setup the dev environment.
.PHONY: git.hooks
git.hooks: tools.verify.go-gitlint
	@-cp -f build/githooks/* .git/hooks/
	@-chmod +x .git/hooks/*

## build: Build source code for host platform.
.PHONY: build
build:
	@$(MAKE) go.build

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

## lint: Check syntax and styling of go sources.
.PHONY: lint
lint:
	@$(MAKE) go.lint

## test: Run unit test.
.PHONY: test
test:
	@$(MAKE) go.test

## cover: Run unit test and get test coverage.
.PHONY: cover
cover:
	@$(MAKE) go.test.cover

## format: Gofmt (reformat) package sources (exclude vendor dir if existed).
.PHONY: format
format:
	@$(MAKE) go.format

## check-updates: Check outdated dependencies of the go projects.
.PHONY: check-updates
check-updates:
	@$(MAKE) go.updates

.PHONY: tidy
tidy:
	@$(GO) mod tidy

## help: Show this help info.
.PHONY: help
help: Makefile
	@echo -e "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"
