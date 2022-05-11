package service

import (
	"TikTok/repository"
	"TikTok/util"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type LoginUser struct {
	User  repository.User
	Token string
}

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

func Login(name string, password string) (int, *LoginUser) {
	code, user := authenticateUser(name, password)
	if code != 0 {
		return code, nil
	}
	token := CreateToken(user)
	if err := RefreshToken(token, user); err != nil {
		return -1, nil
	}
	return 0, &LoginUser{
		User:  *user,
		Token: token,
	}
}

// authenticateUser Returns status code
func authenticateUser(name string, password string) (int, *repository.User) {
	user, err := repository.NewUserDaoInstance().SelectByName(name)
	if err != nil {
		log.Println(err)
		return 1, nil
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return -1, nil
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), hashedPassword)
	if err != nil {
		log.Println(err)
		return 1, nil
	}
	return 0, user
}
