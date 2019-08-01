#!/bin/bash

set -e -u

export PATH=$PATH:$PWD

set -x

export KNI_E2E_NAMESPACE=${1:-kni-test}

go fmt ./tests/...
go test ./tests/ -timeout 60m -test.v $@

echo E2E SUCCESS
