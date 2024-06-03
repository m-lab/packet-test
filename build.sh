#!/bin/sh
set -ex

COMMIT=$(git log -1 --format=%h)
versionflags="-X github.com/m-lab/go/prometheusx.GitShortCommit=${COMMIT}"

go build -v \
    -ldflags "$versionflags" \
     ./cmd/server
go build -v ./cmd/generate-schema
go build -v ./cmd/client/pair1
go build -v ./cmd/client/train1
