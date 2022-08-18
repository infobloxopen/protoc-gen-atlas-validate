include Makefile.buf

GOPATH ?= $(HOME)/go
SRCPATH := $(patsubst %/,%,$(GOPATH))/src

PROJECT_ROOT := github.com/infobloxopen/protoc-gen-atlas-validate

DOCKERFILE_PATH := $(CURDIR)/docker
IMAGE_REGISTRY ?= infoblox
IMAGE_VERSION  ?= dev-atlasvalidate

# configuration for the protobuf gentool
SRCROOT_ON_HOST      := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
SRCROOT_IN_CONTAINER := /home/go/src/$(PROJECT_ROOT)
DOCKERPATH           := /home/go/src
DOCKER_RUNNER        := docker run --rm
DOCKER_RUNNER        += -v $(SRCROOT_ON_HOST):$(SRCROOT_IN_CONTAINER)
DOCKER_GENERATOR     := infoblox/atlas-gentool:$(IMAGE_VERSION)
GENERATOR            := $(DOCKER_RUNNER) $(DOCKER_GENERATOR)

GENVALIDATE_IMAGE      := $(IMAGE_REGISTRY)/atlas-gentool
GENVALIDATE_DOCKERFILE := $(DOCKERFILE_PATH)/Dockerfile

default: vendor options install

.PHONY: vendor
vendor:
	go mod vendor

install:
	go install

regenerate: clean-gen generate

clean-gen:
	cd example/examplepb && rm -f *.pb.atlas.validate.go && rm -f *.pb.gw.go && rm -f *.pb.go
	cd example/external && rm -f *.pb.atlas.validate.go && rm -f *.pb.gw.go && rm -f *.pb.go
	cd options && rm -f *.pb.go

generate: gen-gateway $(BUF) example options

gen-gateway:
	go generate tools/tools.go

.PHONY: example
example:
	buf generate --template example/external/buf.gen.yaml --path example/external
	buf generate --template example/examplepb/buf.gen.yaml --path example/examplepb

.PHONY: options
options:
	buf generate --template options/buf.gen.yaml --path options/atlas_validate.proto

.PHONY: gentool
gentool:
	docker build -f $(GENVALIDATE_DOCKERFILE) -t $(GENVALIDATE_IMAGE):$(IMAGE_VERSION) .
	docker image prune -f --filter label=stage=server-intermediate

gentool-examples: gentool
	$(GENERATOR) \
		-I/go/src/github.com/infobloxopen/protoc-gen-atlas-validate \
		--go_out="$(DOCKERPATH)" --go-grpc_out="$(DOCKERPATH)"\
		--grpc-gateway_out="logtostderr=true:$(DOCKERPATH)" \
		--atlas-validate_out="$(DOCKERPATH)" \
		example/examplepb/example.proto \
		example/examplepb/examplepb.proto \
		example/examplepb/example_multi.proto

	$(GENERATOR) \
		-I/go/src/github.com/infobloxopen/protoc-gen-atlas-validate \
		--go_out="$(DOCKERPATH)" --go-grpc_out="$(DOCKERPATH)"\
		--grpc-gateway_out="logtostderr=true:$(DOCKERPATH)" \
		--atlas-validate_out="$(DOCKERPATH)" \
			example/external/external.proto

gentool-options:
	$(GENERATOR) \
		--go_out="Mgoogle/protobuf/descriptor.proto:$(DOCKERPATH)" \
		$(PROJECT_ROOT)/options/atlas_validate.proto

test: generate
	go test -v -cover ./example/examplepb
