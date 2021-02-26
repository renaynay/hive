# This simulation runs the block production tests.
FROM golang:1-alpine as builder
RUN apk add --update git gcc musl-dev linux-headers

# Build the simulator executable.
ADD . /rpc
WORKDIR /rpc
RUN go build -v .

ENTRYPOINT ["./block-production"]