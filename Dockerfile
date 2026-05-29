# syntax=docker/dockerfile:1

# Step 1: build golang binary
FROM golang:1.17 AS builder
WORKDIR /opt/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-w" -o wechat-webhook

# Step 2: copy binary from step1
FROM alpine:3.18
ENV PATH=/usr/local/bin:$PATH
ENV LANG=C.UTF-8
ENV TZ=Asia/Shanghai

RUN apk add --no-cache ca-certificates openssl wget bash tzdata curl && \
    update-ca-certificates && \
    mkdir -p /usr/bin /usr/sbin /data/wechat-webhook/

COPY --from=builder /opt/app/wechat-webhook /usr/bin/wechat-webhook
COPY start.sh /data/wechat-webhook/
WORKDIR /data/wechat-webhook/
CMD ["/usr/bin/wechat-webhook"]

