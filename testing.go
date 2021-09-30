package testing

import (
	// "crypto/rand"
	// "encoding/base64"
	"fmt"
	"ggapi/db"
	"ggapi/model"
	"reflect"
	//"golang.org/x/crypto/bcrypt"
)

type User model.User
func main(){
	db := db.GetDB()
	a := map[string]interface{}{}
	sql := db.Model(&User{}).Select("id").Where("id = ?",24).First(&a)
	if sql.Error != nil {
		return
	}
	
	
}