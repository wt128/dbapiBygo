package db

import (
	"fmt"
	"os"
	"ggapi/model"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
	err error
)

func Init()  {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}

	dsn := "go_test:" + os.Getenv("DB_KEY") + "@tcp(dockerMySQL)/my_go?charset=utf8mb4&parseTime=True&loc=Local"
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
