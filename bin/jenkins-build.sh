#!/bin/bash

HOME=$(pwd)

echo "export http proxy"
export http_proxy=http://xxxx:1080
export https_proxy=http://xxxx:1080

echo "export GOPATH..."
cd ../../../
path=$(pwd)
export GOPATH=$path
echo "GOPATH:"$GOPATH

echo "build..."
cd $HOME
GIT_COMMIT=$(git rev-parse --short HEAD || echo "GitNotFound")
BUILD_TIME=$(date +%FT%T%z)
echo $GIT_COMMIT $BUILD_TIME
LDFLAGS="-X main.GitCommit=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build -ldflags "${LDFLAGS}" -a -o ${HOME}/server/panda ${HOME}/server/panda.go

if [[ $? -ne 0 ]]; then
    #build error
    echo "build ERROR"
    exit 1
fi
