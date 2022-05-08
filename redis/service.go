package redis

import "TikTok/util"

func Set(key string, value interface{}, expireTime int64) (err error) {
	if expireTime > 0 {
		_, err = Client.Do("SETEX", key, expireTime, value).Result()
	} else {
		_, err = Client.Do("SET", key, value).Result()
	}
	return
}

func Get(key string) (result interface{}, err error) {
	redisVal, err := Client.Do("GET", key).Result()
	if err != nil {
		return nil, err
	}
	err = util.DefaultTranscoder.Unmarshal(redisVal.([]byte), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}