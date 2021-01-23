#DEFAULT MAKEFILE FOR LIBS

.PHONY: test
test: deps
	echo "nothing to test"

.PHONY: build
build: deps
	docker build . -t ${TAG}

.PHONY: deploy
deploy:
	docker push ${TAG}

.PHONY: deps
deps:
	rm -r -f .gateway_api
	mkdir -p .gateway_api
	cp ./../gateway_api/graph/schema.graphqls .gateway_api/schema.graphqls