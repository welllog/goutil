ARG GO_V=1.18
ARG ALPINE_V=3.16

FROM golang:${GO_V} as builder

ENV GO111MODULE=on GOPROXY="https://goproxy.io|https://mirrors.aliyun.com/goproxy/,direct" \
    GOSUMDB="sum.golang.google.cn" \
    GOPRIVATE="gitee.com"

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app

###############################################################################

FROM alpine:$ALPINE_V

LABEL maintainer="orinfy@foxmail.com"

ARG TZ="Asia/Shanghai"

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g'  /etc/apk/repositories \
    && apk update && apk upgrade \
    && apk add tzdata ca-certificates \
    && update-ca-certificates \
    && cp /usr/share/zoneinfo/${TZ} /etc/localtime && echo ${TZ} > /etc/timezone \
    && apk del tzdata \
    # 解析域名优先级
    && echo "hosts: files dns" > /etc/nsswitch.conf

COPY --from=builder /usr/local/bin/app /cmd/app

#EXPOSE 80

VOLUME ["/var/log/go"]

CMD ["/cmd/app"]