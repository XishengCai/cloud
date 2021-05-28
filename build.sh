#!/bin/bash

export GOPROXY=https://goproxy.cn
export GO111MODULE=on

BIN_FILE="cloud"
IMAGE_NAME="cloud"
VERSION="latest"
REGISTRY="xishengcai"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o  ./bin/${BIN_FILE} ./main.go

if [ $? -ne 0 ]; then
    echo "build ERROR"
    exit 1
fi

echo build success

cat <<EOF > Dockerfile
FROM registry.cn-hangzhou.aliyuncs.com/launcher/alpine:latest

MAINTAINER xishengcai <cc710917049@163.com>

COPY ./bin/${BIN_FILE} /usr/local/bin
COPY ./conf  /opt/conf
COPY ./template /opt/template
COPY ./docs /opt/docs

RUN chmod +x /usr/local/bin/${BIN_FILE}

WORKDIR /opt

EXPOSE 80

CMD ["${BIN_FILE}"]
EOF

docker build -t ${REGISTRY}/${IMAGE_NAME}:${VERSION} ./
docker push ${REGISTRY}/${IMAGE_NAME}:${VERSION}
docker rmi ${REGISTRY}/${IMAGE_NAME}:${VERSION}
