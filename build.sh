#!/bin/bash

export GOPROXY=https://goproxy.io

BIN_FILE="cloud"
IMAGE_NAME="cloud"
VERSION="v1"
CURRENT_DIRECTORY=${PWD##*/}

echo "set build info"
GIT_COMMIT=$(git rev-parse --short HEAD || echo "GitNotFound")
BUILD_TIME=$(date +%FT%T%z)

echo $GIT_COMMIT $BUILD_TIME
LDFLAGS="-X main.GitCommit=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o  ./bin/${BIN_FILE} ./main.go

if [ $? -ne 0 ]; then
    echo "build ERROR"
    exit 1
fi
echo build success

cat <<EOF > Dockerfile
FROM alpine:latest

MAINTAINER xishengcai <cc710917049@163.com>

COPY ./bin/${BIN_FILE} /usr/local/bin
COPY ./download /opt/download
COPY ./template /opt/template
COPY ./swaggerui /opt/swaggerui
COPY ./conf   /opt/conf

RUN chmod +x /usr/local/bin/${BIN_FILE}

WORKDIR /opt

EXPOSE 80

CMD ["${BIN_FILE}"]
EOF

docker build -t registry.cn-hongkong.aliyuncs.com/xisheng/${IMAGE_NAME}:${VERSION} ../${CURRENT_DIRECTORY}
docker push registry.cn-hongkong.aliyuncs.com/xisheng/${IMAGE_NAME}:${VERSION}