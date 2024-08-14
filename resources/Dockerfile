# 程序编译镜像
FROM golang:1.22 as build

ENV GO111MODULE=on                               \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /go/release

ADD . .

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ./dist/website-verification ./cmd/website-verification

# 构建镜像
FROM alpine:3.20 as prod

ENV GIN_MODE=release    \
    MQ_HOST="127.0.0.1" \
    MQ_PORT="5672"      \
    MQ_USER="rabbit"    \
    MQ_PWD="123456"     \
    MQ_QUEUE="test"     \
    CONCURRENT=100

# 在build阶段复制时区到
COPY --from=build /go/release/dist/website-verification /website-verification

CMD [ "/website-verification"]