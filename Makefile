# Constants
GHCR_URL=ghcr.io/xplorfin

# variables related to your microservice
# change SERVICE_NAME to whatever you're naming your microservice
SERVICE_NAME=lnurlauth
# name of output binary
BINARY_NAME=lnurlauth

# latest git commit hash
LATEST_COMMIT_HASH=$(shell git rev-parse HEAD)

# go commands and variables
GO=go
GOB=$(GO) build
GOT=$(GO) test
GOM=$(GO) mod

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

macos: check-version gomvendor build-macos

linux: check-version gomvendor build-linux

# Makes sure you're running a version of go which supports
# go modules. 
check-version:
ifeq ($(shell [[ $(MINORVER) -lt 11 ]] && BADVER=yes || BADVER=no; echo $$BADVER), yes)
	@echo "Installed go version ($(GOVERSION)) is lower than 1.11."
	@echo "Go >= 1.11 is required for use with Go modules." 
	@echo "Please update your go version." 
	exit 5
else
ifeq ($(shell [[ $(MINORVER) -lt 14 ]] && LOWVER=yes || LOWVER=no; echo $$LOWVER), yes)
	@echo "Installed go version ($(GOVERSION)) is lower than 1.14."
	@echo "Things will still work, but you should definitely update your installed Go version."
endif
endif

test:
	$(GOT) ./...

gomvendor:
	$(GOM) tidy
	$(GOM) vendor

build-macos:
	env GOOS=$(GOOS_MACOS) GOARCH=$(GOARCH) \
	$(GOB) -mod vendor -o $(BINARY_NAME)

build-linux:
	env GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH) \
	$(GOB) -mod vendor -o $(BINARY_NAME)
 
