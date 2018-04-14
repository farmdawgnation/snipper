#!/bin/bash

GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.BuildVersion=`git describe --tags` -X main.BuildTime=`date -u '+%Y-%m-%d_%I:%M:%S%p'`" -o build/snipper-darwin-amd64 snipper.go
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.BuildVersion=`git describe --tags` -X main.BuildTime=`date -u '+%Y-%m-%d_%I:%M:%S%p'`" -o build/snipper-linux-amd64 snipper.go
