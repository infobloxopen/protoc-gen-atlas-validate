GOPATH ?= $(HOME)/go
SRCPATH := $(patsubst %/,%,$(GOPATH))/src

PROJECT_ROOT := github.com/askurydzin/protoc-gen-atlas-validate

DOCKERFILE_PATH := $(CURDIR)/docker
IMAGE_REGISTRY ?= infoblox
IMAGE_VERSION  ?= dev-atlasvalidate

# configuration for the protobuf gentool
SRCROOT_ON_HOST      := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
SRCROOT_IN_CONTAINER := /go/src/$(PROJECT_ROOT)
DOCKERPATH           := /go/src/
DOCKER_RUNNER        := docker run --rm
DOCKER_RUNNER        += -v $(SRCROOT_ON_HOST):$(SRCROOT_IN_CONTAINER)
DOCKER_GENERATOR     := infoblox/atlas-gentool:$(IMAGE_VERSION)
GENERATOR            := $(DOCKER_RUNNER) $(DOCKER_GENERATOR)

GENVALIDATE_IMAGE      := $(IMAGE_REGISTRY)/atlas-gentool
GENVALIDATE_DOCKERFILE := $(DOCKERFILE_PATH)/Dockerfile

default: vendor options install

.PHONY: vendor
vendor:
	@dep ensure -vendor-only

.PHONY: vendor-update
vendor-update:
	@dep ensure

install:
	go install

.PHONY: gentool
gentool:
	@docker build -f $(GENVALIDATE_DOCKERFILE) -t $(GENVALIDATE_IMAGE):$(IMAGE_VERSION) .
	@docker tag $(GENVALIDATE_IMAGE):$(IMAGE_VERSION) $(GENVALIDATE_IMAGE):latest
	@docker image prune -f --filter label=stage=server-intermediate


gentool-options:
	@$(GENERATOR) \
                --gogo_out="Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:$(DOCKERPATH)" \
                $(PROJECT_ROOT)/options/atlas_validate.proto
