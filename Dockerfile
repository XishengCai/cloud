FROM alpine:latest

COPY ./bin/cloud /usr/local/bin
COPY ./download /opt/download
COPY ./template /opt/template
COPY ./swaggerui /opt/swaggerui
COPY ./conf   /opt/conf

RUN chmod +x /usr/local/bin/cloud

WORKDIR /opt

CMD ["cloud"]
