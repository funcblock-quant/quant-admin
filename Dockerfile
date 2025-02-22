FROM golang:alpine as builder


ENV GOPROXY https://goproxy.cn/

WORKDIR /app
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk add tzdata

COPY go.mod ./
RUN go mod tidy
COPY go.sum ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o quanta-admin .

FROM alpine

WORKDIR /app

COPY --from=builder /app/quanta-admin /app/quanta-admin

EXPOSE 8000

CMD ["/app/quanta-admin","server","-c", "/app/config/settings.${ENVIRONMENT}.yml"]