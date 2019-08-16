from alpine:latest

RUN  mkdir /opt/cloud
ADD ./bin/cloud /opt/cloud
ADD ./conf /opt/cloud

RUN  chmod +x /opt/cloud/cloud

WORKDIR /opt/cloud

cmd ["./cloud"]
