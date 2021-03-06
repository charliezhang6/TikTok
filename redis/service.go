package redis

import (
	"TikTok/util"
	"log"
)

func Set(key string, value interface{}, expireTime int64) (err error) {
	redisValue, _ := util.DefaultTranscoder.Marshal(value)
	if expireTime > 0 {
		_, err = Client.Do("SETEX", key, expireTime, string(redisValue)).Result()
	} else {
		_, err = Client.Do("SET", key, redisValue).Result()
	}
	return
}

func Get(key string, value interface{}) error {
	redisValue, err := Client.Get(key).Result()
	if err != nil {
		log.Println("查询redis失败" + err.Error())
		return err
	}
	err = util.DefaultTranscoder.Unmarshal([]byte(redisValue), &value)
	return err
}
