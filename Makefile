GO = go
GIT = git
DOCKER = docker
INSTALL = sudo install

VERSION ?= 0.0.0
PRERELEASE ?= alpha0

BRANCH ?= $(shell $(GIT) rev-parse --abbrev-ref HEAD 2>/dev/null)
COMMIT ?= $(shell $(GIT) rev-parse HEAD 2>/dev/null)
SHORT_SHA ?= $(shell $(GIT) rev-parse --short $(COMMIT))

REGISTRY ?= ghcr.io
REPOSITORY ?= frantjc/dockerfile-addendum
MODULE ?= github.com/$(REPOSITORY)
TAG ?= latest
IMAGE ?= $(REGISTRY)/$(REPOSITORY):$(TAG)

BUILD_ARGS ?= --build-arg version=$(VERSION) --build-arg prerelease=$(PRERELEASE)

BIN ?= /usr/local/bin

.DEFAULT: install

install: build
	@$(INSTALL) ./bin/addendum $(BIN)

build:
	@$(GO) $@ \
		-ldflags "\
			-s -w \
			-X $(MODULE).Version=$(VERSION) \
			-X $(MODULE).Prerelease=$(PRERELEASE)\
		" \
		-o ./bin \
		./cmd/addendum

image img: 
	@$(DOCKER) build -t $(IMAGE) $(BUILD_ARGS) .

fmt vet test:
	@$(GO) $@ ./...

download vendor verify:
	@$(GO) mod $@

.PHONY: \
	install \
	image img \
	fmt vet test \
	vendor verify
