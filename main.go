package main

import (
	"github.com/joho/godotenv"
	"log"
	"quanta-admin/cmd"
)

//go:generate swag init --parseDependency --parseDepth=6 --instanceName admin -o ./docs/admin

// @title quanta-admin API
// @version 2.0.0
// @description 基于Gin + Vue + Element UI的前后端分离权限管理系统的接口文档
// @description 添加qq群: 521386980 进入技术交流群 请先star，谢谢！
// @license.name MIT
// @license.url https://github.com/go-admin-team/go-admin/blob/master/LICENSE.md

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or could not load it")
	}
	cmd.Execute()
}
