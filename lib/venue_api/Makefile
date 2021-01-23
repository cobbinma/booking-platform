.PHONY: test
test:
	go test -tags integration ./...

.PHONY: build
build: deps
	docker build . -t ${TAG}

.PHONY: deploy
deploy:
	docker push ${TAG}

.PHONY: deps local
local:
	go build cmd/api/main.go
	./main

.PHONY: deps
deps:
	rm -rf .protobuf
	mkdir -p .protobuf
	cp -r ./../protobuf/ .protobuf