#build stage
FROM golang:1.21 AS builder
WORKDIR /usr/src/app
COPY . .
RUN go env -w GO111MODULE=on && \
  go env -w GOPROXY=https://goproxy.cn,direct && \
  go build -o /go/bin/app -v .

#final stage, we need docker cli installed
FROM registry.cn-hangzhou.aliyuncs.com/pivot_studio/pivot_lang:latest
COPY --from=builder /go/bin/app /app
ENTRYPOINT ["/app"]

