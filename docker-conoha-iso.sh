#!/bin/bash

ARGS=$@
NAME=conoha-iso
docker run \
       -ti \
       --rm \
       --name $NAME \
       -e OS_TENANT_NAME=$OS_TENANT_NAME \
       -e OS_PASSWORD=$OS_PASSWORD \
       -e OS_AUTH_URL=$OS_AUTH_URL \
       -e OS_USERNAME=$OS_USERNAME \
       -e OS_REGION_NAME=$OS_REGION_NAME \
       hironobu/conoha-iso \
       /bin/conoha-iso $ARGS
