GIT_SHA1 = $(shell git rev-parse --verify HEAD)
IMAGES_TAG = ${shell git describe --exact-match --tags 2> /dev/null || echo 'latest'}
IMAGE_PREFIX = my-super-awesome-monorepo-

IMAGE_DIRS = $(wildcard lib/* platform/*)

# All targets are `.PHONY` ie always need to be rebuilt
.PHONY: all ${IMAGE_DIRS}

# Build all images
all: ${IMAGE_DIRS}

# Build and tag a single image
${IMAGE_DIRS}:
	$(eval IMAGE_NAME := $(subst /,-,$@))
	echo ${shell pwd}

.PHONY: dev
dev:
	cd lib && docker-compose up --build