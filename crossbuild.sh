#!/bin/bash

GOOS=darwin GOARCH=amd64 go build -o build/snipper-darwin-amd64 cmd/snipper/snipper.go
GOOS=linux GOARCH=amd64 go build -o build/snipper-linux-amd64 cmd/snipper/snipper.go
