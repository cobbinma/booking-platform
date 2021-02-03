.PHONY: test
test:
	echo "nothing to test..."

.PHONY: build
build: deps
	docker build . -t ${TAG}

.PHONY: deploy
deploy:
	docker push ${TAG}

.PHONY: local
local: deps
	cargo run

.PHONY: grpcui
grpcui:
	grpcui -cert localhost.crt -key localhost.key -insecure -import-path .protobuf/ -proto .protobuf/src/booking/api/service.proto localhost:6969

.PHONY: deps
deps:
	rm -rf .protobuf
	mkdir -p .protobuf
	cp -R ./../protobuf/* .protobuf

.PHONY: tests
tests: deps
	cargo test

.PHONY: postgres
postgres:
	docker-compose up