package main

import (
	"go-bot/asynchronous"
	"go-bot/config"
	"go-bot/pkg/memory"
	"go-bot/service"

	"github.com/spf13/pflag"
)

var (
	addr = pflag.StringP("addr", "d", "localhost:8080", "service address")
	cfg  = pflag.StringP("config", "c", "", "apiserver config file path.")
)

func main() {
	pflag.Parse()
	// 开启延时队列

	// 开启日志
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}
	// 开启服务
	memory.InitRedis()

	asynchronous.Run()
	service.LoadService(*addr)
}
