package user

import (
	"github.com/gin-gonic/gin"

	"ggapi/db"
	"ggapi/model"
	_ "net/http"
)

type Service struct{}

type User model.User

func (s Service) CreateModel(c *gin.Context) (User, error) {
	db := db.GetDB()
	var u User

	if err := c.ShouldBindJSON(&u); err != nil {
		return u, err
	}

	if err := db.Create(&u).Error; err != nil {
		return u, err
	}

	return u, nil
}
func (s Service) GetByID(id string) (User, error) {
	db := db.GetDB()
	var u User
	if err := db.Where(" id = ?", id).First(&u).Error; err != nil {
		return u, err
	}
	return u, nil
}

//sessionIdによる検出.
func (s Service) GetByRemember(sid string) (uint, error) {
	db := db.GetDB()
	var u User
	if err := db.Where("remember = ? ", sid).First(&u).Error; err != nil {
		return 0, err
	}
	return u.ID, nil
}

func (s Service) UpdateByID(id string, c *gin.Context) (User, error) {
	db := db.GetDB()
	var u User
	if err := db.Where("id = ?", id).First(&u).Error; err != nil {
		return u, err
	}

	if err := c.BindJSON(&u); err != nil {
		return u, err
	}

	db.Save(&u)

	return u, nil
}

// DeleteByID is delete a User
func (s Service) DeleteByID(id string) error {
	db := db.GetDB()
	var u User

	if err := db.Where("id = ?", id).Delete(&u).Error; err != nil {
		return err
	}

	return nil
}
