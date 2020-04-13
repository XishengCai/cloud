############# builder c #############
FROM golang:1.13.0 as builder

WORKDIR /go/src/cloud
COPY . .
RUN GO111MODULE=off CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install 


############# cloud controller manager #############
FROM alpine:latest
COPY --from=builder /go/bin/cloud /cloud
ENTRYPOINT  ["/cloud"]