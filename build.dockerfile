# This dockerfile uses extends image https://hub.docker.com/woodpecker-kit/woodpecker-npm
# VERSION 1
# Author: woodpecker-kit
# dockerfile offical document https://docs.docker.com/engine/reference/builder/
# https://hub.docker.com/_/golang
FROM golang:1.21.13 as builder

ARG GO_ENV_PACKAGE_NAME=github.com/woodpecker-kit/woodpecker-npm
ARG GO_ENV_ROOT_BUILD_BIN_NAME=woodpecker-npm
ARG GO_ENV_ROOT_BUILD_BIN_PATH=build/${GO_ENV_ROOT_BUILD_BIN_NAME}
ARG GO_ENV_ROOT_BUILD_ENTRANCE=cmd/woodpecker-npm/main.go

ARG GO_PATH_SOURCE_DIR=/go/src
WORKDIR ${GO_PATH_SOURCE_DIR}

RUN mkdir -p ${GO_PATH_SOURCE_DIR}/${GO_ENV_PACKAGE_NAME}
COPY $PWD ${GO_PATH_SOURCE_DIR}/${GO_ENV_PACKAGE_NAME}

# proxy golang
RUN go env -w "GOPROXY=https://goproxy.cn,direct"
RUN go env -w "GOPRIVATE='*.gitlab.com,*.gitee.com"

#RUN go install -v github.com/convention-change/zymosis/cmd/zymosis@v1.1.3
#RUN cd ${GO_PATH_SOURCE_DIR}/${GO_ENV_PACKAGE_NAME} && \
#    zymosis -g go

RUN cd ${GO_PATH_SOURCE_DIR}/${GO_ENV_PACKAGE_NAME} && \
    go mod download -x

RUN  cd ${GO_PATH_SOURCE_DIR}/${GO_ENV_PACKAGE_NAME} && \
  CGO_ENABLED=0 \
  go build \
  -a \
  -installsuffix cgo \
  -ldflags '-w -s --extldflags "-static -fpic"' \
  -tags netgo \
  -o ${GO_ENV_ROOT_BUILD_BIN_PATH} \
  ${GO_ENV_ROOT_BUILD_ENTRANCE}

# https://hub.docker.com/_/node/tags
FROM node:20.11.1-alpine

ARG DOCKER_CLI_VERSION=${DOCKER_CLI_VERSION}
ARG GO_ENV_PACKAGE_NAME=github.com/woodpecker-kit/woodpecker-npm
ARG GO_ENV_ROOT_BUILD_BIN_NAME=woodpecker-npm
ARG GO_ENV_ROOT_BUILD_BIN_PATH=build/${GO_ENV_ROOT_BUILD_BIN_NAME}

ARG GO_PATH_SOURCE_DIR=/go/src

# proxy apk mirrors.aliyun.com
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk --no-cache add \
  git nodejs npm \
  && rm -rf /var/cache/apk/* /tmp/*

RUN mkdir /app
WORKDIR /app

COPY --from=builder ${GO_PATH_SOURCE_DIR}/${GO_ENV_PACKAGE_NAME}/${GO_ENV_ROOT_BUILD_BIN_PATH} .
ENTRYPOINT [ "/app/woodpecker-npm" ]
# CMD ["/app/woodpecker-npm", "--help"]