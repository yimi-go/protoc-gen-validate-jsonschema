SHELL := /bin/bash

# include the common make file
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/../.. && pwd -P))
endif
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/_output
endif
ifeq ($(origin TMP_DIR),undefined)
TMP_DIR := $(OUTPUT_DIR)/tmp
endif

# set the version number. you should not need to do this
# for the majority of scenarios.
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --tags --always --match='v*' 2> /dev/null)
endif
# Check if the tree is dirty.  default to dirty
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2> /dev/null))
	GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD 2> /dev/null)

# Minimum test coverage
ifeq ($(origin COVERAGE),undefined)
COVERAGE := 60
endif

ifeq ($(origin GOOS), undefined)
	GOOS := $(shell go env GOOS)
endif
ifeq ($(origin GOARCH), undefined)
	GOARCH := $(shell go env GOARCH)
endif

FIND := find . ! -name '*.pb.go'
SEDCMD=$(shell which sed)
ifeq ($(shell uname),Darwin)
  SEDCMDI=$(SEDCMD) -i ''
  XARGS := xargs -r
  AWK := gawk
else
  SEDCMDI=$(SEDCMD) -i
  XARGS := xargs --no-run-if-empty
  AWK := awk
endif
GO := go

# Makefile settings
ifndef V
MAKEFLAGS += --no-print-directory
endif

COMMA := ,
SPACE :=
SPACE +=
