FROM alpine
MAINTAINER Hironobu Saitoh(hiro@hironobu.org)

RUN apk update && apk add ca-certificates openssh && rm -rf /var/cache/apk/*
RUN wget https://github.com/hironobu-s/conoha-iso/releases/download/current/conoha-iso-linux.amd64.gz
RUN gunzip -c conoha-iso-linux.amd64.gz > conoha-iso
RUN chmod +x conoha-iso
RUN mv conoha-iso /bin/
