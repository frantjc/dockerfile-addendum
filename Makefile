BIN ?= /usr/local/bin

GO ?= go
GIT ?= git
DOCKER ?= docker

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

INSTALL ?= sudo install

.DEFAULT: install

install: binaries
	@$(INSTALL) $(CURDIR)/bin/sqncd $(CURDIR)/bin/addendum $(BIN)

bins binaries: addendum

addendum:
	@$(GO) build -ldflags "-s -w -X $(MODULE).Version=$(VERSION) -X $(MODULE).Prerelease=$(PRERELEASE)" -o $(CURDIR)/bin $(CURDIR)/cmd/$@

image img: 
	@$(DOCKER) build -t $(IMAGE) $(BUILD_ARGS) .

fmt vet test:
	@$(GO) $@ ./...

tidy vendor verify:
	@$(GO) mod $@

.PHONY: \
	install bins binaries sqnc sqncd \
	shims shimuses shimsource uses source placeholders \
	image img \
	fmt vet test \
	tidy vendor verify \
	clean \
	protos \
	lint tools
