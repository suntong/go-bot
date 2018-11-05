package messages

import (
	"log"

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

var mes *delayqueue

func GetDelayMessages(db string) *delayqueue {
	if mes == nil {
		client := redis.NewClient(&redis.Options{
			Addr: "127.0.0.1:6379",
		})
		_, err := client.Ping().Result()
		if err != nil {
			log.Fatal(err)
		}
		mes = &delayqueue{
			c:  client,
			db: db,
		}
	}

	return mes
}

func GetDefaultMessages() *delayqueue {
	return GetDelayMessages("delaymessages")
}
