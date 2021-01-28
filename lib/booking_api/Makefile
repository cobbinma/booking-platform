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

.PHONY: deps
deps:
	rm -rf .protobuf
	mkdir -p .protobuf
	cp -R ./../protobuf/* .protobuf