GIT_SHA1 = $(shell git rev-parse --verify HEAD)
IMAGES_TAG = ${shell git describe --exact-match --tags 2> /dev/null || echo 'latest'}
IMAGE_PREFIX = booking-

IMAGE_DIRS = $(wildcard lib/*)

.PHONY: all ${IMAGE_DIRS}
all: ${IMAGE_DIRS}

# Build and tag a single image
${IMAGE_DIRS}:
	$(eval IMAGE_NAME := $(subst /,-,$@))
	$(MAKE) -C $@ test
	docker build -t ${DOCKERHUB_OWNER}/${IMAGE_PREFIX}${IMAGE_NAME}:${IMAGES_TAG} -t ${DOCKERHUB_OWNER}/${IMAGE_PREFIX}${IMAGE_NAME}:latest --build-arg TAG=${IMAGE_PREFIX}${IMAGE_NAME} --build-arg GIT_SHA1=${GIT_SHA1} $@
	docker push ${DOCKERHUB_OWNER}/${IMAGE_PREFIX}${IMAGE_NAME}:${IMAGES_TAG}
	docker push ${DOCKERHUB_OWNER}/${IMAGE_PREFIX}${IMAGE_NAME}:latest

.PHONY: tests
tests:
	@for f in ${IMAGE_DIRS}; do $(MAKE) -C $${f} test; done

.PHONY: dev
dev:
	docker-compose up --build