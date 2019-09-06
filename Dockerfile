FROM alpine:latest

MAINTAINER xishengcai <cc710917049@163.com>

RUN apk --no-cache add tzdata zeromq \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo '$TZ' > /etc/timezone

RUN  mkdir /opt/cloud
COPY ./bin/linux/cloud /opt/cloud
COPY ./conf /opt/cloud/

RUN  chmod +x /opt/cloud/cloud

WORKDIR /opt/cloud

CMD ["./cloud", "./conf/cloud_config.toml"]
