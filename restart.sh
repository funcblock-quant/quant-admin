#!/bin/bash
echo "go build"
go mod tidy
go build -o quanta-admin main.go
chmod +x ./quanta-admin
echo "kill quanta-admin service"
killall quanta-admin # kill go-admin service
nohup ./quanta-admin server -c=config/settings.dev.yml >> access.log 2>&1 & #后台启动服务将日志写入access.log文件
echo "run quanta-admin success"
ps -aux | grep quanta-admin
