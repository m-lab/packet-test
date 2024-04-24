#!/bin/sh
set -ex

go build -v ./cmd/server

go build -v ./cmd/generate-schema
