package memory

type livequeue struct {
	db string
}

func (m *livequeue) Push(data interface{}) (int64, error) {
	return client.SAdd(m.db, data).Result()
}

func (m *livequeue) Range() ([]string, error) {
	return client.SMembers(m.db).Result()
}

func GetLive(name string) *livequeue {
	return &livequeue{
		db: name,
	}
}

// type liveKV struct {
// 	db string
// }

// func (m *liveKV) Get(data string) (string, error) {
// 	return client.HGet(m.db, data).Result()
// }

// func (m *liveKV) Range() ([]string, error) {
// 	return client.LRange(m.db, 0, -1).Result()
// }

// func GetLiveKV(name string) *liveKV {
// 	return &liveKV{
// 		db: name,
// 	}
// }

type kv struct {
	db string
}

func (m *kv) Get(key string) (string, error) {
	return client.HGet(m.db, key).Result()
}

func (m *kv) GetKey() ([]string, error) {
	return client.HKeys(m.db).Result()
}

func (m *kv) Del(key ...string) (int64, error) {
	return client.HDel(m.db, key...).Result()
}

func (m *kv) Set(key string, value string) (bool, error) {
	return client.HSet(m.db, key, value).Result()
}

func GetKV(name string) *kv {
	return &kv{
		db: name,
	}
}
