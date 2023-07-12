# syntax=docker/dockerfile:1
FROM golang:1.20.5-bullseye AS protoc

RUN apt-get update && apt-get install -y \
      protobuf-compiler=3.12.4-1 \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /work

COPY go.mod go.mod
COPY go.sum go.sum

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
