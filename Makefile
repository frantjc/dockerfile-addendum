GO = go
GOLANGCI-LINT = golangci-lint
DOCKER = docker
INSTALL = sudo install

VERSION ?= 0.0.0
PRERELEASE ?= alpha0

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

lint:
	@$(GOLANGCI-LINT) run --fix

.PHONY: \
	install \
	image img \
	fmt vet test \
	download vendor verify \
	lint
