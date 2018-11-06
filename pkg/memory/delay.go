package memory

import (
	"github.com/lexkong/log"
	"github.com/spf13/viper"

	"github.com/go-redis/redis"
)

type delayqueue struct {
	c  *redis.Client
	db string
}

func (d *delayqueue) Push(data interface{}) (int64, error) {
	return d.c.RPush(d.db, data).Result()
}

func (d *delayqueue) Pop() (string, error) {
	resultSlice, err := d.c.BLPop(0, d.db).Result()
	if err == nil && resultSlice[0] == d.db {
		return resultSlice[1], nil
	}
	return "", err
}

func (d *delayqueue) Close() error {
	return d.c.Close()
}

var DefaultMes *delayqueue

func InitRedis() {
	if DefaultMes == nil {
		DefaultMes = GetDefaultMessages()
	}
}

func GetDelayMessages(db string) *delayqueue {
	client := redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.addr"),
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Error("delay", err)
	}
	tmp := &delayqueue{
		c:  client,
		db: db,
	}
	return tmp
}

func GetDefaultMessages() *delayqueue {
	return GetDelayMessages(viper.GetString("redis.delay_name"))
}
