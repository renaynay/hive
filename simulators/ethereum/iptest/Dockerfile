# This simulation runs various graphql tests
FROM golang:1.13-alpine as builder
RUN apk add --no-cache git make gcc musl-dev linux-headers bash  curl jq
RUN GOPATH=/go

# Add the local test files
ADD /ethereum/graphql/graphql.go /go/src/github.com/ethereum/hive/simulators/ethereum/iptest/iptest.go
ADD common /go/src/github.com/ethereum/hive/simulators/common

RUN cd /go/src/github.com/ethereum/hive/simulators/ethereum/iptest/ && go build iptest.go && go run iptest.go
