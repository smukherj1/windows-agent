#!/usr/bin/bash

set -eux

PROTOC_OPTS="--go_out=grpc \
--go_opt=module=github.com/smukherj1/windows-agent/grpc \
--go-grpc_out=grpc --go-grpc_opt=module=github.com/smukherj1/windows-agent/grpc"

protoc ${PROTOC_OPTS} proto/server.proto

go build -o out/server bin/server/server.go

GOOS=windows GOARCH=amd64 go build -o out/agent.exe bin/agent/agent.go
go build -o out/agent bin/agent/agent.go
