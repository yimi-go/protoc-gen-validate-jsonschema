GOBIN ?=$(shell go env GOBIN)
ifeq ($(GOBIN),)
	GOBIN=$(HOME)/go/bin
endif

default: generate

.PHONY: install
install:
	@echo "Installing protoc-gen-jsonschema-test plugin for test bin file generating..."
	@go install ../../../cmd/protoc-gen-jsonschema-test

.PHONY: generate
generate: install
	@echo "Generating protoc generate request bin files for testing..."
	@find . -name "*.proto" | while read pbfile; do buf generate --path "$$pbfile"; done

.PHONY: uninstall
uninstall:
	@echo "Uninstalling protoc-gen-jsonschema-test plugin..."
	@rm -f $(GOBIN)/protoc-gen-jsonschema-test
