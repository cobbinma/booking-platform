.PHONY: test
test:
	go test -tags integration ./...

.PHONY: build
build:
	docker build . -t ${TAG}

.PHONY: deploy
deploy:
	docker push ${TAG}

.PHONY: generate
generate:
	gqlgen generate
	go generate ./...

.PHONY: local
local:
	go build cmd/api/main.go
	./main