DOCKER:=$(shell which docker)
ifeq ($(DOCKER),)
$(error "Docker is not installed. Please install Docker to use this Makefile.")
endif

REGISTRY?=docker.io
ORG?=${USER}
REPO_PREFIX?=""
REPO?=$(shell basename $(shell pwd))
IMAGE?=$(REGISTRY)/$(ORG)/$(REPO_PREFIX)${REPO}

.PHONY: all clean squeaky
all: no-embed embed

embed: .embed
.embed: Dockerfile.embedded sphinx.tmpl .no-embed
	${DOCKER} build -t $(IMAGE):embed -f $< . && touch $@

no-embed: .no-embed
.no-embed: Dockerfile
	${DOCKER} build -t $(IMAGE):no-embed -f $< . && touch $@

clean:
	@rm -f .no-embed .embed

squeaky: clean
	@echo "Cleaning up images..."
	@${DOCKER} rmi -f $(IMAGE):no-embed $(IMAGE):embed || true

push: .no-embed
	@echo "Pushing images..."
	@${DOCKER} push $(IMAGE):no-embed
	@${DOCKER} push $(IMAGE):embed