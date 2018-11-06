package memory

import (
	"github.com/go-redis/redis"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

var client *redis.Client

func setClient() {
	client = redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.addr"),
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Error("client", err)
	}
}

func InitRedis() {
	setClient()
	InitMes()
}
