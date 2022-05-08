package service

import (
	"TikTok/repository"
	"TikTok/util"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func RegisterUser(name string, password string) (int64, string, error) {
	token := name + password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	userId := util.GenSonyflake()
	user := repository.User{ID: userId, Name: name, Password: string(hashedPassword)}
	err = repository.NewUserDaoInstance().AddUser(user)
	if err != nil {
		log.Fatal(err)
		return 0, "", err
	}
	return userId, token, nil
}
