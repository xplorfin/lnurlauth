# Constants
GHCR_URL=ghcr.io/xplorfin

# variables related to your microservice
# change SERVICE_NAME to whatever you're naming your microservice
SERVICE_NAME=lnurlauth
# name of output binary
BINARY_NAME=main

# latest git commit hash
LATEST_COMMIT_HASH=$(shell git rev-parse HEAD)

# go commands and variables
GO=go
GOB=$(GO) build
GOM=$(GO) mod

# git commands
GIT=git

# environment variables related to
# cross-compilation.
GOOS_MACOS=darwin
GOOS_LINUX=linux
GOARCH=amd64

# currently installed/running Go version (full and minor)
GOVERSION=$(shell go version | grep -Eo '[1-2]\.[[:digit:]]{1,3}\.[[:digit:]]{0,3}')
MINORVER=$(shell echo $(GOVERSION) | awk '{ split($$0, array, ".") } {print array[2]}')

# Color code definitions
# Note: everything is bold.
GREEN=\033[1;38;5;70m
BLUE=\033[1;38;5;27m
LIGHT_BLUE=\033[1;38;5;32m
MAGENTA=\033[1;38;5;128m
RESET_COLOR=\033[0m

COLORECHO = $(1)$(2)$(RESET_COLOR)

default: help

# build a vendored binary for macos
macos: gomvendor build-macos

# build a vendored binary for linx
linux: gomvendor build-linux

setup-hooks: ## setup the repository (enables git hooks)
	git config core.hooksPath .github/hooks --replace-all

bench:  ## Run all benchmarks in the Go application
	@go test -bench=. -benchmem

clean-mods: ## Remove all the Go mod cache
	@go clean -modcache

coverage: ## Get the test coverage from go-coverage
	@go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

godocs: ## Run a godoc server
	@echo "godoc server running on http://localhost:9000"
	@godoc -http=":9000"


golangci-install:
	@#Travis (has sudo)
	@if [ "$(shell which golangci-lint)" = "" ] && [ $(TRAVIS) ]; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s && sudo cp ./bin/golangci-lint $(go env GOPATH)/bin/; fi;
	@#AWS CodePipeline
	@if [ "$(shell which golangci-lint)" = "" ] && [ "$(CODEBUILD_BUILD_ID)" != "" ]; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin; fi;
	@#Github Actions
	@if [ "$(shell which golangci-lint)" = "" ] && [ "$(GITHUB_WORKFLOW)" != "" ]; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s && sudo cp ./bin/golangci-lint $(go env GOPATH)/bin/; fi;
	@#Brew - MacOS
	@if [ "$(shell which golangci-lint)" = "" ] && [ "$(shell which brew)" != "" ]; then brew install golangci-lint; fi;

go-acc-install:
	@if [ "$(shell which "go-acc")" = "" ]; then go get -u github.com/ory/go-acc; fi;

lint: golangci-install ## Run the golangci-lint application (install if not found) & fix issues if possible
	@golangci-lint run --fix

# pre-commit hook
pre-commit: lint

test: ## run tests without coverage reporting
	@go test ./...

ci-test: go-acc-install # run a test with coverage
	@go-acc -o profile.cov ./...

gomvendor: ## run tidy & vendor
	@go mod tidy
	@go mod vendor

build-macos: ## build binaries for macos
	env GOOS=$(GOOS_MACOS) GOARCH=$(GOARCH) \
	$(GOB) -mod vendor -o $(BINARY_NAME)

build-linux: ## build binaries for linux
	env GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH) \
	$(GOB) -mod vendor -o $(BINARY_NAME)


help: ## This help dialog.
	@IFS=$$'\n' ; \
	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//'`); \
	for help_line in $${help_lines[@]}; do \
		IFS=$$'#' ; \
		help_split=($$help_line) ; \
		help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		printf "%-30s %s\n" $$help_command $$help_info ; \
	done

docker: docker-build docker-push ## build and push docker image to GHCR

# Builds a docker container from a Dockerfile.
# If you have more complex Docker arguments than just tags and the source,
# consider writing your own recipe.
DOCKER_IMAGE_NAME=$(GHCR_URL)/$(SERVICE_NAME)
docker-build: docker/Dockerfile.make gomvendor ## build a docker file
	@echo "[*] Latest git commit hash: $(call COLORECHO,$(GREEN),$(LATEST_COMMIT_HASH))"
	@echo "[*] Building Docker image $(call COLORECHO,$(BLUE),$(DOCKER_IMAGE_NAME))" \
	"with tags $(call COLORECHO,$(MAGENTA),latest), $(call COLORECHO,$(LIGHT_BLUE),$(LATEST_COMMIT_HASH))"
	docker build -f docker/Dockerfile.make \
		-t $(DOCKER_IMAGE_NAME):latest \
		-t $(DOCKER_IMAGE_NAME):$(LATEST_COMMIT_HASH) \
		.

docker-push: ## pushes a built image to GHCR
	@echo "[*] Pushing $(call COLORECHO,$(LIGHT_BLUE),$(DOCKER_IMAGE_NAME):$(LATEST_COMMIT_HASH)) to GHCR"
	docker push $(DOCKER_IMAGE_NAME):$(LATEST_COMMIT_HASH)
	@echo "[*] Pushing $(call COLORECHO,$(MAGENTA),$(DOCKER_IMAGE_NAME):latest) to GHCR"
	docker push $(DOCKER_IMAGE_NAME):latest