package repository

import (
	"log"
	"sync"
)

type User struct {
	ID       int64  `gorm:"column:user_id"`
	Name     string `gorm:"column:user_name"`
	Password string `gorm:"column:user_password"`
	Salt     string `gorm:":column:salt"`
	Token    string `gorm:"column:token"`
}

type UserDao struct {
}

var userDao *UserDao
var userOnce sync.Once

func NewUserDaoInstance() *UserDao {
	userOnce.Do(
		func() {
			userDao = &UserDao{}
		})
	return userDao
}

func (*UserDao) AddUser(user User) error {
	err := db.Create(&user).Error
	if err != nil {
		log.Fatal("用户添加失败")
		return err
	}
	return nil
}
