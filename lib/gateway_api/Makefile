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
	go get github.com/99designs/gqlgen
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
	cp -R ./../protobuf/* .protobuf

.PHONY: coverage
coverage:
	go test -v -coverpkg=./... -coverprofile=profile.cov ./... && go tool cover -html=profile.cov