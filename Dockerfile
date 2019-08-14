from alpine:latest
ADD ./bin/cloud /opt/

RUN  chmod +x /opt/cloud

WORKDIR /opt/

cmd [".cloud"]
