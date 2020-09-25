# syntax = docker/dockerfile:experimental
FROM golang:1.15-buster AS builder

ENV GO111MODULE=on
WORKDIR /src

COPY . .

RUN go test ./... \
  && go test -tags integration ./... \
  && CGO_ENABLED=0 go build -a -o /main cmd/api/*.go

# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
FROM alpine:3.12.0 as release

RUN apk add --update \
  curl \
  tini \
  ;

COPY --from=builder /main /

EXPOSE 8989
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/main"]
