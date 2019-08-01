#!/bin/bash

set -e -x -u

BUILD_VALUES= ./hack/build.sh

GOOS=darwin GOARCH=amd64 go build -o kni-darwin-amd64 ./cmd/...
GOOS=linux GOARCH=amd64 go build -o kni-linux-amd64 ./cmd/...
GOOS=windows GOARCH=amd64 go build -o kni-windows-amd64.exe ./cmd/...

shasum -a 256 ./kni-*-amd64*
