#!/bin/bash

path=$(pwd)
export GOPATH=$path
echo "GOPATH:"$GOPATH

echo "build..."
GIT_COMMIT=$(git rev-parse --short HEAD || echo "GitNotFound")
BUILD_TIME=$(date +%FT%T%z)

echo $GIT_COMMIT $BUILD_TIME
LDFLAGS="-X main.GitCommit=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}"

CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
go build -ldflags "${LDFLAGS}" -a -o ${HOME}/cloud ${HOME}/cloud/main.go

if [[ $? -ne 0 ]]; then
    #build error
    echo "build ERROR"
    exit 1
fi