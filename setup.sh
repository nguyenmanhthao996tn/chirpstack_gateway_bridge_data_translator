#!/bin/bash

# https://grpc.io/docs/languages/go/quickstart/

go mod tidy

# install the protocol compiler plugins for Go
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
go get github.com/chirpstack/chirpstack/api/go/v4
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc


# protoc --go_out=. *.proto

# go run .