APP := $(notdir $(shell pwd))
IMAGE_NAME := $(APP)
BUILD_DATE ?= $(shell date +"%Y%m%dT%H%M%SZ")

GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT = $(shell git rev-list --tags --max-count=1)
SEMVER := $(shell git describe --tags $(GIT_COMMIT))
VERSION := $(SEMVER)-$(BUILD_DATE).$(GIT_BRANCH)

IMPORT_PATH := github.com/porter614/$(APP)
ARTIFACTORY_USER ?= $(shell whoami)

DEV_DOCKER_REPO ?= artifactory.com:8157

GOPROXY_ARG ?= $(GOPROXY)

all: install

install:
	go install -v -ldflags "-X $(IMPORT_PATH)/pkg/version.App=$(APP) \
		-X $(IMPORT_PATH)/pkg/version.Version=$(VERSION) \
		-X $(IMPORT_PATH)/pkg/version.SemVer=$(SEMVER) \
		-X $(IMPORT_PATH)/pkg/version.BuildDate=$(BUILD_DATE) \
		-X $(IMPORT_PATH)/pkg/version.BuildUser=$(ARTIFACTORY_USER) \
		-X $(IMPORT_PATH)/pkg/version.CommitId=$(GIT_COMMIT) \
		-X $(IMPORT_PATH)/pkg/version.Branch=$(GIT_BRANCH) \
		-X $(IMPORT_PATH)/cmd.App=$(APP)"

build:
	docker build \
		--build-arg APP=$(APP) \
		--build-arg GOPROXY_ARG=$(GOPROXY_ARG) \
		--target build -t $(APP)_build .

clean:
	rm -f ./$(APP)

distclean: clean
	rm -f $(GOPATH)/bin/$(APP)
	# Note that this could fail if you have a ton of dangly things.  Some advanced docker wrangling may be required.
	# Order is to avoid removing layers with dependents
	docker images -a | grep "$(APP)_cross-compile" | awk '{print $$3}' | xargs docker rmi
	docker images -a | grep "$(APP)_test" | awk '{print $$3}' | xargs docker rmi
	docker images -a | grep "$(APP)_lint" | awk '{print $$3}' | xargs docker rmi
	docker images -a | grep "$(APP)_build" | awk '{print $$3}' | xargs docker rmi

lint:
	docker build \
		--build-arg APP=$(APP) \
		--build-arg GOPROXY_ARG=$(GOPROXY_ARG) \
		--target lint -t $(APP)_lint .

test:
	docker build --rm=false \
		--build-arg APP=$(APP) \
		--build-arg GOPROXY_ARG=$(GOPROXY_ARG) \
		--build-arg BUILD_NUMBER=$(BUILD_NUMBER) \
		--target test -t $(APP)_test .

image:
	docker build --rm=false \
		--build-arg APP=$(APP)\
		--build-arg GOPROXY_ARG=$(GOPROXY_ARG) \
		--build-arg VERSION=$(VERSION) \
		--build-arg SEMVER=$(SEMVER) \
		--build-arg BUILD_USER=$(ARTIFACTORY_USER) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg GIT_BRANCH=$(GIT_BRANCH) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_NUMBER=$(BUILD_NUMBER) \
		--build-arg IMPORT_PATH=$(IMPORT_PATH) \
		--label "semver=$(SEMVER)" \
		--label "version=$(VERSION)" \
		--label "build_date_time=$(BUILD_DATE)" \
		--label "branch=$(GIT_BRANCH)" \
		--label "commit_id=$(GIT_COMMIT)" \
		--label "build_user=$(ARTIFACTORY_USER)" -t $(APP) .

eng: VERSION := $(SEMVER)-ENG-$(ARTIFACTORY_USER)
eng: image
	docker tag $(APP):latest \
		$(DEV_DOCKER_REPO)/$(APP):$(VERSION)
ifdef PUSH
	@docker login \
		-u $(ARTIFACTORY_USER) \
		-p $(ARTIFACTORY_API_KEY) \
		$(DEV_DOCKER_REPO)
	docker push $(DEV_DOCKER_REPO)/$(APP):$(VERSION)
endif

version:
	@echo $(VERSION)

build-native:
	go build -ldflags "-X $(IMPORT_PATH)/pkg/version.App=$(APP) \
		-X $(IMPORT_PATH)/pkg/version.Version=$(VERSION) \
		-X $(IMPORT_PATH)/pkg/version.SemVer=$(SEMVER) \
		-X $(IMPORT_PATH)/pkg/version.BuildDate=$(BUILD_DATE) \
		-X $(IMPORT_PATH)/pkg/version.BuildUser=$(ARTIFACTORY_USER) \
		-X $(IMPORT_PATH)/pkg/version.CommitId=$(GIT_COMMIT) \
		-X $(IMPORT_PATH)/pkg/version.Branch=$(GIT_BRANCH)" 

.PHONY: install build lint test clean distclean image eng version



