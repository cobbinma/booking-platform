.PHONY: test
test:
	go test -tags integration ./...

.PHONY: build
build: deps
	docker build . -t ${TAG}

.PHONY: deploy
deploy:
	docker push ${TAG}

.PHONY: generate
generate:
	gqlgen generate
	go generate ./...

.PHONY: local
local: deps
	go build cmd/api/main.go
	./main

.PHONY: deps
deps:
	rm -rf .protobuf
	mkdir -p .protobuf
	cp -r ./../protobuf/ .protobuf
	ls -la