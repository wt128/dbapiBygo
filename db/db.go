package db

import (
	"ggapi/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
	err error
)

func Init()  {
	dsn := "root:kfs985A&@tcp(127.0.0.1:3306)/mygo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn),&gorm.Config{})
	if err != nil {
		panic(err)
	}
	autoMigrate()
}

func GetDB() *gorm.DB {
	return db
}

func autoMigrate() {
	db.AutoMigrate(&model.Post{})
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Like{})
	db.AutoMigrate(&model.CntList{})
}
