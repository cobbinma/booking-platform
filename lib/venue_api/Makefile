.PHONY: test
test:
	go test -tags integration ./...

.PHONY: build
build:
	docker build . -t ${TAG}

.PHONY: deploy
deploy:
	docker push ${TAG}