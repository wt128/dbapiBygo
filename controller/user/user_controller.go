package user

import (
	"fmt"
	"ggapi/db"
	"ggapi/model"
	"ggapi/service"
	"net/http"
	"strings"
	"github.com/go-playground/validator"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Controller struct{}
type User model.User
type Post model.Post
//user_id int

func (pc Controller) Index(c *gin.Context){
	var p []Post
	db := db.GetDB()
	sql := db.Where("user_id = ?",c.Params.ByName("id")).Find(&p)
	if sql.Error != nil{
		c.AbortWithStatus(400)
		return
	}
	c.JSON(200,p)
	}

func (pc Controller) Avatar(c *gin.Context){
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatus(403)
	}
	c.SaveUploadedFile(file,"../../view")
	c.JSON()
}
// 07 /  13 validationをなんとかする
func (pc Controller) Create(c *gin.Context){
	var u User
	db := db.GetDB()
	validate := validator.New()
	
	if err := c.ShouldBindJSON(&u); err != nil {
		
		c.JSON(http.StatusBadRequest, gin.H{"errorMessages": err.Error()})
		c.Abort()
		return
	}

	if err := validate.Struct(&u); err != nil{
		var errorMessages []string 
		for _, err := range err.(validator.ValidationErrors) {
			var errorMessage string
			fieldName := err.Field() //バリデーションでNGになった変数名を取得
		
		switch fieldName {
		
			case "Email":
			  var typ = err.Tag() //バリデーションでNGになったタグ名を取得
			  switch typ {
				case "required":
					errorMessage = "ここは必須です"
				case "email":
					errorMessage = "メアドが不適切です。"
				}
			case "Password":
				errorMessage = "パスワード入れてね"
			}
			errorMessages = append(errorMessages, errorMessage)
		  }	
		c.JSON(403, gin.H{"errorMessages": errorMessages})
		
		//c.Redirect(302,"/user")
		return
	}
	
	passwrd := []byte(u.Password)
	hashed, _ := bcrypt.GenerateFromPassword(passwrd,9)
	u.Password = string(hashed)
	strings.ToLower(u.Email)
	u.Remember = model.UsernewToken()
    
	if err := db.Create(&u).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessages": "このメールアドレスはすでに使用されています"})
		fmt.Println(err)
		return 
	}

	c.SetCookie("sid",u.Remember,3600*24*30*12,"/","localhost:3000",true,false)
	c.JSON(200,gin.H{"sid":u.Remember})
}

func (pc Controller) Show(c *gin.Context){
	id := c.Params.ByName("id")
	var posts[] Post
	
	db := db.GetDB()
	if err := db.Where("userid = ?",id).Find(&posts).Error; err != nil{
		c.AbortWithStatus(403)
		fmt.Println(err)
		return
	}
	c.JSON(204,gin.H{})
	
}

// Delete action: DELETE /users/:id
func (pc Controller) Delete(c *gin.Context) {
	type Sid struct{
		Id	string	`json:"sid"`
	}

    id := c.Params.ByName("id")
	db := db.GetDB()
    var s user.Service
	var sid Sid
		if err := c.BindJSON(&sid); err != nil {
		c.AbortWithStatus(400)
	}
	
	checkSid := db.Model(&User{}).Where("remember = ?",sid.Id).First(&User{})
	if checkSid.Error != nil {
		c.AbortWithStatusJSON(403,gin.H{"errorMessages":"不正なリクエストです。[session error]"})
		return
	}

    if err := s.DeleteByID(id); err != nil {
        c.AbortWithStatus(403)
        fmt.Println(err)
    } else {
        c.JSON(204, gin.H{"id #" + id: "deleted"})
    }
}

func (pc Controller) Update(c *gin.Context){
	type New struct{ //requestのjson形式
		ID		    uint	`json:"id"`
		Email	    string	`json:"email" validate:"email"`
		Confirm     string	`json:"confirm"`
		Password	string	`json:"password" validate:"required"`
	}
	var u New
	validate := validator.New()
	var prepass string
	db := db.GetDB()
	
	if err := c.BindJSON(&u); err != nil {
		c.AbortWithStatus(400)
	}

	passConfirm := db.Model(&User{}).Where("id = ?",u.ID).Pluck("password",&prepass)
	if passConfirm.Error != nil{
		c.AbortWithStatus(403)
		
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(prepass),[]byte(u.Confirm)); err!= nil{
		fmt.Println("ログinできません")
		c.AbortWithStatusJSON(403, gin.H{"errorMessages": "元のパスワードが間違ってますよ。"})
		c.Abort()
		return
	}
	
	if err := validate.Struct(&u); err != nil{
		var errorMessages []string 
		for _, err := range err.(validator.ValidationErrors) {
			var errorMessage string
			fieldName := err.Field() //バリデーションでNGになった変数名を取得
		
		switch fieldName {
			case "Email":
				var typ = err.Tag() //バリデーションでNGになったタグ名を取得
				switch typ {
				case "required":
					errorMessage = "ここは必須です"
				case "email":
					errorMessage = "メアドが不適切です。"
				}
			case "Password":
				errorMessage = "パスワード入れてね"
			}
			errorMessages = append(errorMessages, errorMessage)
			}	
		c.AbortWithStatusJSON(403, gin.H{"errorMessages": errorMessages})
		c.Abort()
		return
	}
	passwrd := []byte(u.Password)
	hashed, _ := bcrypt.GenerateFromPassword(passwrd,9)
	u.Password = string(hashed)
	strings.ToLower(u.Email)

	sql := db.Model(&User{}).Where("id = ?",u.ID).Updates(User{Email:u.Email, Password:string(hashed)})
	
	if sql.Error != nil {
		c.AbortWithStatusJSON(403,gin.H{"errorMessages":"エラーが発生しました(入力されたアドレスはもうつかわれてるかも？)"})
		return
	}

	c.JSON(204,gin.H{"ok":"成功しました。"})
	

}
