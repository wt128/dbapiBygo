package session

import (
	"fmt"
	"ggapi/db"
	"ggapi/model"
	
	"strings"

	// "ggapi/service"
	"net/http"

	//"github.com/go-playground/validator"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User model.User
type Post model.Post
type Sid struct{ ID string `json:"sid"` }

type Controller struct{}

func (pc Controller) Create(c *gin.Context){
	var u User
	if err := c.BindJSON(&u); err != nil{
		c.AbortWithStatus(400)
		return
	}
	setUser, err := getUser(strings.ToLower(u.Email))
	if err == false {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessages": "メアドが不適切です。"})
		 //どうにかしとく
		c.Abort()
		return
	}
	dbPass := []byte(setUser.Password)
	formPass := []byte(u.Password)
	//fmt.Println(setUser.Remember)
	if err := bcrypt.CompareHashAndPassword(dbPass,formPass); err!= nil{
		fmt.Println("ログinできません")
		c.JSON(http.StatusBadRequest, gin.H{"errorMessages": "パスワードが間違ってますよ。"})
		c.Abort()
		} else{
			
			fmt.Println("ログインできました。")
			c.SetCookie("sid",setUser.Remember,3600*24*365*10,"/","localhost",true,true)
			c.JSON(200, gin.H{"sid": setUser.Remember}) //どうにかしとく
		}
}

func (pc Controller)SetId(c *gin.Context){

	var s Sid
	result := map[string]interface{}{}
	db := db.GetDB()
	
	if err := c.BindJSON(&s); err != nil{
		c.AbortWithStatus(400)
		return
	}

	err := db.Model(&User{}).Select("id").Where("remember = ?",s.ID).First(&result)
	if err.Error != nil {
		c.AbortWithStatus(400)
		return
	}
	fmt.Println(result)
	c.JSON(200,result)
	
}

func (pc Controller)GetEmail (c *gin.Context){
	var sid Sid
	result := map[string]interface{}{}
	db := db.GetDB()
	if err := c.BindJSON(&sid); err != nil{
		c.AbortWithStatusJSON(403,gin.H{"errorMessages":"不正なリクエストです。"})
	}
	sql := db.Model(&User{}).Select("email").Where("remember = ?",sid.ID).Scan(&result)
	if sql.Error != nil{
		c.AbortWithStatusJSON(403,gin.H{"errorMessages":"エラーが発生しました"})
	}
	c.JSON(200,result)
}

func getUser(email string) (User, bool){
	db := db.GetDB()
	var user User
	result := db.First(&user,"email = ?",email)
	if result.Error != nil{
		return user,false
	}
	return user,true
}
