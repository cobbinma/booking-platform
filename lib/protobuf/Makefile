.PHONY: test
test:
	echo "nothing to test"

.PHONY: build
build: gen-lang-go

.PHONY: deploy
deploy:
	echo "nothing to deploy"

SORTED_PROTO_FILES=$(shell find src -name '*.proto' | sort)

gen-lang-go:
	go mod download
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	rm -rf autogen/lang/go
	mkdir -p autogen/lang/go
	for PROTOFILE in ${SORTED_PROTO_FILES}; do \
		protoc -I ./ --go_out=plugins=grpc:${GOPATH}/src $$PROTOFILE \
	; done