#!/bin/bash

export GOPROXY=https://goproxy.io

echo "set build info"
GIT_COMMIT=$(git rev-parse --short HEAD || echo "GitNotFound")
BUILD_TIME=$(date +%FT%T%z)

echo $GIT_COMMIT $BUILD_TIME
LDFLAGS="-X main.GitCommit=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ./bin/linux/cloud

if [ $? -ne 0 ]; then
    echo "build failed"
else
    echo "build succeed"
fi
