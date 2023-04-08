#!/bin/bash

# https://go.dev/doc/install/source#environment

read -p 'Target OS: ' TARGET_OS
read -p 'Target Architecture: ' TARGET_ARCH

GOOS=$TARGET_OS GOARCH=$TARGET_ARCH go build -ldflags "-s -w" plow/zapper

# Run garble too?