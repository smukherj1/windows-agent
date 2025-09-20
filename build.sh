#!/usr/bin/bash

set -eux

go build -o out/server bin/server/server.go

GOOS=windows GOARCH=amd64 go build -o out/agent.exe bin/agent/agent.go
go build -o out/agent bin/agent/agent.go