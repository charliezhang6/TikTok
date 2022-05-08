package redis

func Set(key string, value interface{}, expireTime int64) (err error) {
	if expireTime > 0 {
		_, err = Client.Do("SETEX", key, expireTime, value).Result()
	} else {
		_, err = Client.Do("SET", key, value).Result()
	}
	return
}

func Get(key string) (interface{}, error) {
	redisVal, err := Client.Do("GET", key).Result()
	if err != nil {
		return nil, err
	}
	return redisVal, nil
}
