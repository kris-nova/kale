TARGET = kale
GOTARGET = github.com/kris-nova/$(TARGET)
REGISTRY ?= krisnova
IMAGE = $(REGISTRY)/$(TARGET)
DIR := ${CURDIR}
DOCKER ?= docker

GIT_VERSION ?= $(shell git describe --always --dirty)
IMAGE_VERSION ?= $(shell git describe --always --dirty)
IMAGE_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD | sed 's/\///g')
GIT_REF = $(shell git rev-parse --short=8 --verify HEAD)

default: compile

all: compile install

push: ## Push to the docker registry
	$(DOCKER) push $(REGISTRY)/$(TARGET):$(GIT_REF)

clean: ## Clean the docker images
	rm -f $(TARGET)
	$(DOCKER) rmi $(REGISTRY)/$(TARGET) || true

container: ## Build the docker container
	$(DOCKER) build \
		-t $(REGISTRY)/$(TARGET):$(IMAGE_VERSION) \
		-t $(REGISTRY)/$(TARGET):$(IMAGE_BRANCH) \
		-t $(REGISTRY)/$(TARGET):$(GIT_REF) \
	    -t $(REGISTRY)/$(TARGET):latest \
		.

run: ## Run the controller in a container
	$(DOCKER) run $(REGISTRY)/$(TARGET):$(IMAGE_VERSION)


compile: ## Compile the binary into bin/kale
	go build -o bin/kale main.go

install: ## Create the kale executable in $GOPATH/bin directory.
	install -m 0755 bin/kale ${GOPATH}/bin/kale

.PHONY: help
help:  ## Show help messages for make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'
