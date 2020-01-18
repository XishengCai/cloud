FROM alpine:latest

MAINTAINER xishengcai <cc710917049@163.com>

COPY ./bin/cloud /usr/local/bin
COPY ./download /opt/download
COPY ./template /opt/template
COPY ./swaggerui /opt/swaggerui
COPY ./conf   /opt/conf

RUN chmod +x /usr/local/bin/cloud

WORKDIR /opt

CMD ["cloud"]
