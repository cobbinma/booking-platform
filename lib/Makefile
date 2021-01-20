#DEFAULT MAKEFILE FOR LIBS

.PHONY: test
test:
	echo "nothing to test"

.PHONY: build
build:
	docker build . -t ${TAG}

.PHONY: deploy
deploy:
	docker push ${TAG}