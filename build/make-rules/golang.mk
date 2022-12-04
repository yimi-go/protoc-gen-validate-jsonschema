GO_LDFLAGS += -X $(VERSION_PACKAGE).GitVersion=$(VERSION) \
	-X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
ifneq ($(DLV),)
	GO_BUILD_FLAGS += -gcflags "all=-N -l"
	LDFLAGS = ""
endif
GO_BUILD_FLAGS += -ldflags "$(GO_LDFLAGS)"

ifeq ($(GOOS),windows)
	GO_OUT_EXT := .exe
endif

ifeq ($(ROOT_PACKAGE),)
	$(error the variable ROOT_PACKAGE must be set prior to including golang.mk)
endif

GOPATH := $(shell $(GO) env GOPATH)
ifeq ($(origin GOBIN), undefined)
	GOBIN := $(GOPATH)/bin
endif

COMMANDS ?= $(shell $(FIND) -name .cmd -type f | sed -e "s|/.cmd$$||g")
BINS ?= $(foreach cmd,${COMMANDS},$(notdir ${cmd}))
ifeq (${BINS},)
  $(error Could not determine BINS)
endif

EXCLUDE_TESTS=$(shell cat $(ROOT_DIR)/.gotestignore)

.PHONY: go.build.%
go.build.%:
	$(eval COMMAND := $*)
	@echo "===========> Building binary $(COMMAND) $(VERSION) for $(GOOS) $(GOARCH)"
	@mkdir -p $(OUTPUT_DIR)/platforms/$(GOOS)/$(GOARCH)
	@CGO_ENABLED=0 $(GO) build $(GO_BUILD_FLAGS) -o $(OUTPUT_DIR)/platforms/$(GOOS)/$(GOARCH)/$(COMMAND)$(GO_OUT_EXT) $(ROOT_PACKAGE)/cmd/$(COMMAND)

.PHONY: go.build
go.build: $(addprefix go.build., $(BINS))

.PHONY: go.clean
go.clean:
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

.PHONY: go.lint
go.lint: tools.verify.golangci-lint
	@echo "===========> Run golangci to lint source codes"
	@golangci-lint run -c $(ROOT_DIR)/.golangci.yaml $(ROOT_DIR)/...

.PHONY: go.format
go.format: tools.verify.golines tools.verify.goimports
	@echo "===========> Formating codes"
	@$(FIND) -type f -name '*.go' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' | $(XARGS) goimports -w -local $(ROOT_PACKAGE)
	@$(FIND) -type f -name '*.go' | $(XARGS) golines -w --max-len=120 --reformat-tags --shorten-comments --ignore-generated .
	@go mod edit -fmt

ifeq ($(origin V), 1)
GO_TEST_FLAG += -v
endif

.PHONY: go.test
go.test: tools.verify.go-junit-report
	@mkdir -p $(OUTPUT_DIR)
	@echo "===========> Run unit test"
	@$(GO) test $(GO_BUILD_FLAGS) $(GO_TEST_FLAG) $(ROOT_PACKAGE)/...
	@set -o pipefail;$(GO) test -race -cover -coverprofile=$(OUTPUT_DIR)/coverage.out \
		-timeout=10m -shuffle=on -short $(GO_TEST_FLAG) `go list ./...|\
		egrep -v $(subst $(SPACE),'|',$(sort $(EXCLUDE_TESTS)))` 2>&1 | \
		tee >(go-junit-report --set-exit-code >$(OUTPUT_DIR)/report.xml)
	@$(SEDCMDI) '/mock_.*.go/d' $(OUTPUT_DIR)/coverage.out # remove mock_.*.go files from test coverage
	@$(SEDCMDI) '/.*.pb.go/d' $(OUTPUT_DIR)/coverage.out # remove .*.pb.go files from test coverage
	@$(GO) tool cover -html=$(OUTPUT_DIR)/coverage.out -o $(OUTPUT_DIR)/coverage.html

.PHONY: go.test.cover
go.test.cover: go.test
	@$(GO) tool cover -func=$(OUTPUT_DIR)/coverage.out | \
		$(AWK) -v target=$(COVERAGE) -f $(ROOT_DIR)/build/coverage.awk

.PHONY: go.updates
go.updates: tools.verify.go-mod-outdated
	@$(GO) list -u -m -json all | go-mod-outdated -update -direct
