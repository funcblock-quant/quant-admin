FROM golang:alpine as builder


ENV GOPROXY https://goproxy.cn/

WORKDIR /app
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk add tzdata

COPY go.mod go.sum ./
RUN go mod tidy # 整理依赖
RUN go mod download # 下载依赖
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o quanta-admin .

FROM alpine

WORKDIR /app

COPY --from=builder /app/quanta-admin /app/quanta-admin
COPY --from=builder /app/config/settings.yml /app/config/settings.yml
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

EXPOSE 8000

CMD ["/app/quanta-admin","server","-c", "/app/config/settings.yml"]