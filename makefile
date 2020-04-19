SHA1         		:= $(shell git rev-parse --verify --short HEAD)
MAJOR_VERSION		:= $(shell cat app.json | sed -n 's/.*"major": "\(.*\)",/\1/p')
MINOR_VERSION		:= $(shell cat app.json | sed -n 's/.*"minor": "\(.*\)"/\1/p')
INTERNAL_BUILD_ID	:= $(shell [ -z "${TRAVIS_BUILD_NUMBER}" ] && echo "local" || echo ${TRAVIS_BUILD_NUMBER})
BINARY				:= $(shell cat app.json | sed -n 's/.*"name": "\(.*\)",/\1/p')
VERSION				:= $(shell echo "${MAJOR_VERSION}.${MINOR_VERSION}.${INTERNAL_BUILD_ID}-${SHA1}")
BUILD_IMAGE			:= $(shell echo "golang:1.14.2")
PWD					:= $(shell pwd)

ENV 				?= local

.DEFAULT_GOAL := build

.PHONY: version
version:
	@echo "Setting build to Version: v$(VERSION)" 
	$(shell echo v$(VERSION) > VERSION.txt)

.PHONY: fmt
fmt: ## Runs go fmt on code base
ifeq ($(ENV),local)
	go fmt ./...
else
	docker run --rm --name=$(BUILD_NAME) \
	-v $(PWD):/usr/src/$(REPO) \
	-w /usr/src/$(REPO) $(BUILD_IMAGE) \
	go fmt ./...
endif

.PHONY: lint
lint: ## Runs more than 20 different linters using golangci-lint to ensure consistency in code.
ifeq ($(ENV),local)
	golangci-lint run -v
else
	docker run --rm --name=$(BUILD_NAME) \
	-e GOPACKAGESPRINTGOLISTERRORS=1 \
	-e GO111MODULE=on \
	-e REPO=$(REPOPATH)$(REPO) \
	-v $(PWD):/usr/src/ \
	golangci/golangci-lint:v1.24 \
	/bin/bash -c 'cd /usr/src/ && golangci-lint run'
endif

.PHONY: test
test: ## Runs the tests within a docker container
ifeq ($(ENV),local)
	go test -cover -v -p 8 -count=1 ./...
else
	docker run --rm --name=$(BUILD_NAME) \
	-v $(PWD):/usr/src/$(REPO) \
	-w /usr/src/$(REPO) $(BUILD_IMAGE) \
	go test -cover -v -p 8 -count=1 ./...
endif

.PHONY: build
build: version test
	@echo "Building"

	docker run -it --rm \
	-v $(PWD):/usr/src/myapp \
	-w /usr/src/myapp $(BUILD_IMAGE) \
	bash scripts/docker/build.sh $(BINARY) $(VERSION)
	
	@echo "Completed build"

.PHONY: release
release:
	@echo "Releasing"

	docker run -it --rm \
	-v $(PWD):/usr/src/myapp \
	-w /usr/src/myapp $(BUILD_IMAGE) \
	bash scripts/docker/release.sh $(VERSION)
	
	@echo "Completed release"