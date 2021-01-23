.PHONY: test
test:
	echo "nothing to test"

.PHONY: build
build:
	echo "nothing to build"

.PHONY: deploy
deploy:
	echo "nothing to deploy"

SORTED_PROTO_FILES=$(shell find src -name '*.proto' | sort)

gen-lang-go:
	go mod download
	go install github.com/golang/protobuf/protoc-gen-go
	rm -rf autogen/lang/go
	mkdir -p autogen/lang/go
	for PROTOFILE in ${SORTED_PROTO_FILES}; do \
		protoc $$PROTOFILE -I. --go_out=plugins=grpc:${GOPATH}/src \
	; done