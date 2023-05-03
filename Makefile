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

define copy_installer
	cp scripts/install.sh ${RELEASE_DIR}/
	chmod u+x ${RELEASE_DIR}/install.sh
endef

define go_build
	go build -o ${RELEASE_BASE} -ldflags "-X ${MODULE}/pkg/version.Version=$(1)"
endef

define go_release_build
	GOOS=$(2) GOARCH=$(3) go build -o ${RELEASE_BASE}-$(2)-$(3)$(4) -ldflags "-X ${MODULE}/pkg/version.Version=$(1)"
endef

######## Goals

build:
	@echo -n 'Building PowerSlide CLI... '
	@$(call make_dir,${RELEASE_DIR})
	@$(call copy_installer,${RELEASE_DIR})
	@$(call go_build,dev.${DATE})
	@echo 'Done!'

release:
	@echo -n 'Building PowerSlide CLI ${RELEASE_VERSION}... '
	@$(call make_dir,${RELEASE_DIR})
	@$(call copy_installer,${RELEASE_DIR})
	@$(call go_build,${RELEASE_VERSION})

	@$(call go_release_build,${RELEASE_VERSION},linux,386)
	@$(call go_release_build,${RELEASE_VERSION},linux,amd64)
	@$(call go_release_build,${RELEASE_VERSION},linux,arm)
	@$(call go_release_build,${RELEASE_VERSION},linux,arm64)

	@$(call go_release_build,${RELEASE_VERSION},darwin,amd64)
	@$(call go_release_build,${RELEASE_VERSION},darwin,arm64)

	@$(call go_release_build,${RELEASE_VERSION},windows,386,.exe)
	@$(call go_release_build,${RELEASE_VERSION},windows,amd64,.exe)
	@$(call go_release_build,${RELEASE_VERSION},windows,arm,.exe)
	@$(call go_release_build,${RELEASE_VERSION},windows,arm64,.exe)

	@sha256sum ${RELEASE_BASE}-* > ${RELEASE_DIR}/sha256sums
	@echo 'Done!'

test_release:
	@echo -n 'Building PowerSlide CLI ${RELEASE_VERSION}... '
	@$(call make_dir,${RELEASE_DIR})
	@$(call copy_installer,${RELEASE_DIR})
	@$(call go_build,${RELEASE_VERSION})
	@echo 'Done!'

clean:
	@rm -rf ${RELEASE_DIR}
	@echo "Project cleaned"

changelog:
	@scripts/generate-changelog.sh
