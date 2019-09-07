FROM alpine:latest

MAINTAINER xishengcai <cc710917049@163.com>

RUN apk --no-cache add tzdata zeromq \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo '$TZ' > /etc/timezone

RUN  mkdir /opt/cloud
COPY . /opt/cloud

WORKDIR /opt/cloud

RUN  chmod +x /opt/cloud/bin/linux/cloud && \
     mkdir /etc/cloud && \
     cp /opt/cloud/bin/linux/cloud /usr/local/bin && \
     cp ./conf/cloud_config.toml /etc/cloud/

CMD ["cloud", "/etc/cloud/cloud_config.toml"]
