FROM golang:1.8

MAINTAINER yumaito

RUN cp -p /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

RUN apt-get update && \
    apt-get install -y redis-tools

ENV project_repo github.com/yumaito/darumasan
ENV app_path /go/src/${project_repo}
RUN mkdir -p ${app_path}
WORKDIR ${app_path}

RUN curl https://glide.sh/get | sh

RUN go get -u -v gopkg.in/godo.v2/cmd/godo
RUN mkdir -p ${app_path}/Gododir
