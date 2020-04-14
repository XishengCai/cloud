#!/bin/bash

set -e

BIN_FILE="cloud"
IMAGE_NAME="cloud"
VERSION="latest"

# echo "set build info"
# GIT_COMMIT=$(git rev-parse --short HEAD || echo "GitNotFound")
# BUILD_TIME=$(date +%FT%T%z)

# echo $GIT_COMMIT $BUILD_TIME
# LDFLAGS="-X main.GitCommit=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}"

# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o  ./bin/${BIN_FILE} ./main.go

# if [ $? -ne 0 ]; then
#     echo "build ERROR"
#     exit 1
# fi
# echo build success

cat <<EOF > Dockerfile
############# builder c #############
FROM golang:1.13.0 as builder

WORKDIR /go/src/cloud
COPY . .
RUN GO111MODULE=off CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install 


############# cloud controller manager #############
FROM alpine:latest
COPY --from=builder /go/bin/cloud /cloud
ENTRYPOINT  ["/cloud"]
EOF

docker build -t xishengcai/${IMAGE_NAME}:${VERSION} ./
docker push xishengcai/${IMAGE_NAME}:${VERSION}
docker rmi xishengcai/${IMAGE_NAME}:${VERSION}




#########  kaniko  #############

# docker run -ti --rm -v `pwd`:/workspace \
# -v /root/.docker/config.json:/kaniko/.docker/config.json:ro gcr.io/kaniko-project/executor:latest \
# --dockerfile=Dockerfile \
# --destination=registry.cn-hongkong.aliyuncs.com/launcher/cloud:latest