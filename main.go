package main

import (
	"flag"
	"go-bot/service"
)

var addr = flag.String("addr", "localhost:8080", "service address")

func main() {
	flag.Parse()
	// 开启延时队列

	// 开启数据库

	// 开启日志

	// 开启服务
	service.LoadService(*addr)
}
