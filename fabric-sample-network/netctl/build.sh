#!/bin/bash

set -e

# build netctl binary
CGO_ENABLED=0 go build

# build Docker image netctl:latest
docker build -t netctl:latest .
