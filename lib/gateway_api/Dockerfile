# syntax = docker/dockerfile:experimental
FROM golang:1.15-buster AS builder

ENV GO111MODULE=on
WORKDIR /src

COPY .protobuf ./.protobuf

COPY go.mod .
COPY go.sum .

RUN ls
RUN ./.protobuf

RUN go mod download

COPY . .

RUN go test ./... \
  && CGO_ENABLED=0 go build -a -o /main cmd/api/*.go

# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
FROM alpine:3.12.0 as release

RUN apk add --update \
  curl \
  tini \
  ;

COPY --from=builder /main /

EXPOSE 9999
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/main"]
