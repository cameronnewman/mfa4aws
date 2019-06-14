SHA1         		:= $(shell git rev-parse --verify --short HEAD)
MAJOR_VERSION		:= $(shell cat app.json | sed -n 's/.*"major": "\(.*\)",/\1/p')
MINOR_VERSION		:= $(shell cat app.json | sed -n 's/.*"minor": "\(.*\)"/\1/p')
INTERNAL_BUILD_ID	:= $(shell [ -z "${TRAVIS_BUILD_NUMBER}" ] && echo "local" || echo ${TRAVIS_BUILD_NUMBER})
BINARY				:= $(shell cat app.json | sed -n 's/.*"name": "\(.*\)",/\1/p')
VERSION				:= $(shell echo "${MAJOR_VERSION}.${MINOR_VERSION}.${INTERNAL_BUILD_ID}-${SHA1}")
BUILD_IMAGE			:= $(shell echo "golang:1.12.5")
PWD					:= $(shell pwd)

.DEFAULT_GOAL := build

.PHONY: version
version:
	@echo "Setting build to Version: v$(VERSION)" 
	$(shell echo v$(VERSION) > VERSION.txt)

.PHONY: test
test: version
	@echo "Running Tests"

	docker run -e -it --rm \
	-v $(PWD):/usr/src/myapp \
	-w /usr/src/myapp $(BUILD_IMAGE) \
	bash -c "go test -cover -race -coverprofile=coverage.txt -covermode=atomic -v ./... -count=1"
	@echo "Completed tests"

.PHONY: build
build: test
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