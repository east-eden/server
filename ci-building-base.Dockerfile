# build server base image

FROM golang:1.16.5-alpine3.13

RUN apk add build-base
RUN apk add openssh
RUN apk add make
RUN apk add git
RUN apk add --update docker openrc
RUN rc-update add docker boot