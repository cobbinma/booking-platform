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
	rm -rf autogen/lang/go
	mkdir -p autogen/lang/go
	for PROTOFILE in ${SORTED_PROTO_FILES}; do \
		protoc --proto_path=./ --go_out=autogen/lang/go --go_opt=paths=source_relative $$PROTOFILE \
	; done