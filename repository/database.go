package repository

import (
	"TikTok/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() error {
	var err error
	dsn := config.Config.MySQL.User + ":" + config.Config.MySQL.Password + "@tcp(" + config.Config.MySQL.Host + ":" + config.Config.MySQL.Port + ")/" + config.Config.MySQL.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}
