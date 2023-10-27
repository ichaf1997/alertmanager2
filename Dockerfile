FROM golang:1.20.10-alpine3.18 as builder

ARG GO111MODULE=on
ARG GOPROXY=https://goproxy.io

WORKDIR /opt/src/pcg/alertmanager2

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk add --no-cache make \
    && rm -rf /var/cache/apk/*

COPY . /opt/src/pcg/alertmanager2

RUN make

#----------------------------------------------------------------

FROM alpine:3.18.4

LABEL .image.authors="ichaff1997@gmail.com"

ARG USER="monitor"

WORKDIR /app

COPY --from=builder /opt/src/pcg/alertmanager2/bin/alertmanager2 /app/alertmanager2

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk add --no-cache tzdata\
    && ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && rm -rf /var/cache/apk/* \
    && adduser -D $USER \
    && chown $USER:$USER /app/alertmanager2 \
    && chmod 700 /app/alertmanager2

COPY --chown=$USER:$USER templates/*.tmpl /app/templates/

USER $USER

ENTRYPOINT ["/app/alertmanager2"]

CMD ["--web.listen-address=:8080"]