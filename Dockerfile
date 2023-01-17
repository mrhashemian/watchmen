# 1. Stage for building the app for developement usage (Docker compose or simple docker build command)
FROM golang:1.16-alpine AS build

RUN mkdir -p /src

WORKDIR /src

COPY go.mod go.sum /src/
RUN go mod download

COPY . /src
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-w -s"

# 2. Stage for running the app build in stage 1 for using in developement (docker compose or simple docker build)
FROM alpine:3.11

ENV TZ=Asia/Tehran \
    PATH="/app:${PATH}"

ARG BUILD_PATH

RUN apk add --update tzdata ca-certificates bash && \
    cp --remove-destination /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone && \
    mkdir -p /app && \
    mkdir -p /var/log && \
    chgrp -R 0 /var/log && \
    chmod -R g=u /var/log && \
    chgrp -R 0 /app && \
    chmod -R g=u /app

WORKDIR /app

COPY $BUILD_PATH /app

CMD ["./watchmen", "serve", "--config", "config.yml"]
