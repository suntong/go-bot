package memory

import (
	"fmt"

	"github.com/spf13/viper"
)

type delayqueue struct {
	db string
}

func (d *delayqueue) Push(data interface{}) (int64, error) {
	fmt.Println(data)
	value, ok := data.([]interface{})
	if ok {
		var err error
		for i := range value {
			_, err = client.RPush(d.db, value[i]).Result()
		}
		if err != nil {
			return 0, err
		}
		return 1, nil
	}
	return client.RPush(d.db, data).Result()
}

func (d *delayqueue) Pop() (string, error) {
	resultSlice, err := client.BLPop(0, d.db).Result()
	if err == nil && resultSlice[0] == d.db {
		return resultSlice[1], nil
	}
	return "", err
}

func (d *delayqueue) Close() error {
	return client.Close()
}

var DefaultMes *delayqueue

func InitMes() {
	if DefaultMes == nil {
		DefaultMes = GetDefaultMessages()
	}
}

func GetDelayMessages(db string) *delayqueue {

	tmp := &delayqueue{
		db: db,
	}

	return tmp
}

func GetDefaultMessages() *delayqueue {
	return GetDelayMessages(viper.GetString("redis.delay_name"))
}
