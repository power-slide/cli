# Default goal
.DEFAULT_GOAL := build

######## Constants

override RELEASE_DIR = bin
override RELEASE_BASE = ${RELEASE_DIR}/pwrsl
override MODULE = github.com/power-slide/cli
override DATE := $(shell date +%Y-%m-%d.%H:%M:%S)

######## Variables

ifndef RELEASE_VERSION
override RELEASE_VERSION = test-release.${DATE}
endif

######## Functions

define make_dir
	mkdir -p $(1)
endef

define go_build
	go build -o ${RELEASE_BASE} -ldflags "-X ${MODULE}/pkg/version.Version=$(1)"
endef

######## Goals

build:
	@echo -n 'Building PowerSlide CLI... '
	@$(call make_dir,${RELEASE_DIR})
	@$(call go_build,dev.${DATE})
	@echo 'Done!'

clean:
	@rm -rf ${RELEASE_DIR}
	@echo "Project cleaned"
