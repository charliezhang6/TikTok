package service

import (
	"TikTok/config"
	"TikTok/redis"
	"TikTok/repository"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"log"
)

func CreateToken() string {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal("生成token出错" + err.Error())
	}
	text := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(text))
	token := string(hash[:])
	return token
}

func RefreshToken(token string, user *repository.User) error {
	err := redis.Set(config.LoginTokenKey+token, user, config.Config.ExpireTime*60)
	return err
}
