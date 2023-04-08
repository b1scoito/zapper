#!/bin/bash

# https://go.dev/doc/install/source#environment
# Strip debug info and run garble
GOOS=linux GOARCH=amd64 garble -literals -tiny -seed=random build -ldflags "-s -w" plow/zapper