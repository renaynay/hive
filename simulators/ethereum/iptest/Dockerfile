# This simulation runs various graphql tests
FROM golang:1.13-alpine as builder
RUN apk add --no-cache git make gcc musl-dev linux-headers bash  curl jq
RUN GOPATH=/go

ENV GO111MODULE=off
RUN go get -u -v github.com/derekparker/delve/cmd/dlv
ENV GO111MODULE=on

# Add the local test files
ADD /ethereum/graphql/graphql.go /go/src/github.com/ethereum/hive/simulators/ethereum/iptest/iptest.go
ADD common /go/src/github.com/ethereum/hive/simulators/common

RUN cd /go/src/github.com/ethereum/hive/simulators/ethereum/iptest/
RUN go build
RUN go run iptest.go
