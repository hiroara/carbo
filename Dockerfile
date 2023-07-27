# syntax=docker/dockerfile:1
FROM golang:1.20.5-bullseye AS base

WORKDIR /usr/local/go/src/github.com/hiroara/carbo


FROM base AS protoc

RUN apt-get update && apt-get install -y \
      protobuf-compiler=3.12.4-1 \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

COPY go.mod go.mod
COPY go.sum go.sum

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0


FROM base AS godoc

RUN go install golang.org/x/tools/cmd/godoc@v0.11.0

EXPOSE 6060

CMD ["godoc", "-http=:6060"]


FROM base AS formatter

RUN go install mvdan.cc/gofumpt@v0.5.0
