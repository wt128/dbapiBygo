package post

import (
	"fmt"
	"os"
	"ggapi/db"
	"ggapi/model"
	"math"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	imgupload "github.com/olahol/go-imageupload"
)

type User model.User
type Post model.Post
type Controller struct{}

func (pc Controller) Show(c *gin.Context){
	var post []Post
	db := db.GetDB()
	page, _ := strconv.Atoi(c.Query("p"))
	if page == 0 { page = 1 }
	var offset int = 8 * (page - 1)
	
	sql := db.Offset(offset).Limit(8).Find(&post)
	if sql.Error != nil {
		c.JSON(403,gin.H{"errorMessages":"エラーが発生しました。"})
		c.Abort()
	}

	c.JSON(200,post)
}

func (pc Controller) PostTotal(c *gin.Context){
	var a []int //一応ページ総数
	b := make([]int, 1,1)
	db := db.GetDB()
	sql := db.Model(&Post{}).Select("COUNT(ID)").Pluck("COUNT(ID)",&a)

	if sql.Error != nil {
		c.AbortWithStatus(403)
	}

	b[0] = int(math.Ceil( float64(a[0]) / 8.0))

	c.JSON(200,b)
}

// Likeテーブルとの関連付け
// 
func (pc Controller) Post(c *gin.Context){
	var post Post
	db := db.GetDB()
	if err := c.BindJSON(&post); err != nil{
		c.AbortWithStatus(403)
	}

	if err := db.Create(&post).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessages": "エラーが発生しました。"})
		fmt.Println(err)
		return
	}
	
	c.String(200,strconv.Itoa(int(post.ID)))

}

func (pc Controller) Setimg(c *gin.Context){
	var id []uint
	db := db.GetDB()
	pid := c.Params.ByName("pid")

	file, err := imgupload.Process(c.Request,"file")
	sid := c.Request.FormValue("sid")
	fmt.Println(sid)

	if err != nil {
		c.AbortWithStatus(403)
	}

	checkSid := db.Model(&User{}).Where("remember = ?",sid).Pluck("id",&id)
	if checkSid.Error != nil {
		c.AbortWithStatusJSON(403,gin.H{"errrorMessages":"不正なリクエストです。"})
		c.Abort()
	}

	thumb, err := imgupload.ThumbnailJPEG(file, 300, 300, 90)
	if err != nil {
		panic(err)
	}
	sql := db.Model(&Post{}).Where("id = ?",pid).Update("isimg",true)
	if sql.Error != nil {
		c.AbortWithStatus(403)
	}
	fmt.Println(id)
	thumb.Save(fmt.Sprintf("./view/posts/%d%s.jpg",id[0],pid))
	//c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file))
	c.String(http.StatusAccepted,"ok")
}

func (pc Controller) Getimg(c *gin.Context){
	var uid []uint
	db := db.GetDB()
	pid := c.Params.ByName("pid")
	if err := db.Model(&Post{}).Where("id = ?",pid).Pluck("user_id",&uid).Error; err!= nil{
		c.AbortWithStatus(403)
	}
	c.File(fmt.Sprintf("view/posts/%d%s.jpg",uid[0],pid))
}

func (pc Controller) GetOne(c *gin.Context){
	var post Post
	db := db.GetDB()
	pid := c.Params.ByName("id")
	if err := db.Where("ID = ?",pid).First(&post).Error; err != nil{
		c.JSON(403,gin.H{"errorMessage":"読み込みに失敗しました。"})
	}
	
	c.JSON(200,post)
}

func (pc Controller) Delete(c *gin.Context){
	type DeleteOne struct {
		PostId uint `json:"postid"`
		Sid string `json:"sid"`
	}
	var p DeleteOne
	db := db.GetDB()

	if err := c.BindJSON(&p); err != nil{
		c.AbortWithStatus(400)
		return
	}
	uid := checkSid(p.Sid)

	if uid == 0 {
		c.AbortWithStatusJSON(403,gin.H{"errorMessages":"不正なリクエストです."})
		return
	}

	if err := db.Delete(&Post{},p.PostId).Error; err != nil{
		c.AbortWithStatusJSON(403,gin.H{"errorMessages":"エラーが発生しました"})
		return
	}
	os.Remove(fmt.Sprintf("view/posts/%d%d.jpg",uid,p.PostId))
	c.JSON(200,gin.H{
		"succcess":"削除しました.","uid":uid})
}

func (pc Controller) Update(c *gin.Context){
	type UpdateOne struct {
		Id 		uint 	`json:"id"`
		Sid 	string 	`json:"sid"`
		Content string 	`json:"content"`
		Title	string 	`json:"title"`
	}
	var p UpdateOne
	db := db.GetDB()
	
	if err := c.BindJSON(&p); err != nil{
		c.AbortWithStatus(400)
	}

	if checkSid(p.Sid) == 0{
		c.AbortWithStatusJSON(403,gin.H{"errorMessages":"不正なリクエストです。 [session Error]"})
		return
	}

	sql := db.Model(&Post{}).Where("id = ?",p.Id).Updates(Post{Title:p.Title, Content:p.Content})
	if sql.Error != nil {
		c.AbortWithStatusJSON(403,gin.H{"errorMessages":"エラーが発生しました."})
		return
	}

	c.String(200, strconv.Itoa(int(p.Id)))
}


func checkSid(id interface{}) (interface{}){
	db := db.GetDB()
	uid := map[string]interface{}{}

	sql := db.Model(&User{}).Select("id").Where("remember = ?",id).Find(&uid)
	result := uid["id"]
	if err := sql.Error; err != nil{
		return 0
	}
	return result
}