.PHONY: test
test:
	go test -tags integration ./...

.PHONY: build
build: deps
	docker build . -t ${TAG}

.PHONY: deploy
deploy:
	docker push ${TAG}

.PHONY: local
local: deps
	go build cmd/api/main.go
	./main

.PHONY: deps
deps:
	rm -rf .protobuf
	mkdir -p .protobuf
	cp -R ./../protobuf/* .protobuf

.PHONY: dev
dev: deps
	docker-compose up --build

.PHONY: grpcui
grpcui:
	grpcui -cert localhost.crt -key localhost.key -insecure -import-path .protobuf/ -proto .protobuf/src/customer/api/service.proto localhost:8899