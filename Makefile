GO = go
GIT = git
GOLANGCI-LINT = golangci-lint

SEMVER ?= 1.0.1

fmt generate test:
	@$(GO) $@ ./...

download vendor verify:
	@$(GO) mod $@

lint:
	@$(GOLANGCI-LINT) run --fix

release:
	@$(GIT) tag -a v$(SEMVER) -m v$(SEMVER)
	@$(GIT) push --follow-tags

gen: generate
dl: download
ven: vendor
ver: verify
format: fmt

.PHONY: up fmt generate test download vendor verify lint shim clean gen dl ven ver format release
