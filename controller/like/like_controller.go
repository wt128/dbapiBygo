package like

import (
	"fmt"
	"ggapi/db"
	"ggapi/model"
	"github.com/gin-gonic/gin"
)

type Like model.Like
type CntList model.CntList
type Controller struct{}

func (pc Controller) Create(c *gin.Context) {
	var com Like
	db := db.GetDB()
	//currrent_user
	if err := c.BindJSON(&com); err != nil {
		c.JSON(400, gin.H{"errorMessages": "エラーが発生しました."})
		return
	}
	if err := db.Create(&com).Error; err != nil {
		c.JSON(400, gin.H{"errorMessages": "DBとの接続がうまくできませんでした."})
		return
	}
}

// @like = Like.find_by(post_id: params[:post_id], user_id: current_user.id)
// @like.destroy
func (pc Controller) Destroy(c *gin.Context) {
	var com Like
	db := db.GetDB()

	if err := c.BindJSON(&com); err != nil {
		c.JSON(400, gin.H{"errorMessages": "エラーが発生しました."})
		return
	}
	err := db.Where("post_id = ? and user_id= ?", com.PostID, com.UserID).Delete(&com)
	if err.Error != nil {
		c.JSON(400, gin.H{"errorMessages": err.Error})
		return
	}
	c.JSON(204, gin.H{"success": "できました。"})
}

func (pc Controller) Show(c *gin.Context) {
	var cnt []CntList

	db := db.GetDB()
	err := db.Model(&Like{}).Select("COUNT(id) as cnt", "post_id as pid").Group("pid").Find(&cnt)
	fmt.Println(cnt)

	if err.Error != nil {
		c.JSON(400, gin.H{"errorMessages": "エラーが発生しました。"})
		return
	}
	c.JSON(200, cnt)
}

func (pc Controller) Auser(c *gin.Context) {
	var cnt []CntList
	id := c.Params.ByName("id")
	db := db.GetDB()

	sql := db.Model(&Like{}).Where("user_id = ? ", id).Select("COUNT(id) as cnt", "post_id as pid").Group("pid").Find(&cnt)
	if sql.Error != nil {
		c.AbortWithStatus(403)
		fmt.Println(sql.Error)
		return
	}
	c.JSON(200, cnt)
}

func (pc Controller) Check(c *gin.Context) {
	result := map[uint]interface{}{}

	db := db.GetDB()
	sql := db.Model(&Like{}).Select("post_id").Where("user_id = ?", c.Params.ByName("id")).Find(&result)
	if sql.Error != nil {
		c.AbortWithStatus(400)
		return
	}
	fmt.Println(result)
	c.JSON(200, result)
}

//db.Table("(?) as u", db.Model(&User{}).Select("name", "age")).Where("age = ?", 18}).Find(&User{})
// SELECT * FROM (SELECT `name`,`age` FROM `users`) as u WHERE `age` = 18

//select * from  (SELECT COUNT(Likes.id) as cnt,Likes.PostId as pid FROM `cnt_lists` GROUP BY `PostId`)
// b.Select("AVG(age) as avgage").Group("name").Having("AVG(age) > (?)", subQuery).Find(&results)
// // SELECT AVG(age) as avgage FROM `users` GROUP BY `name` HAVING AVG(age) > (SELECT AVG(age) FROM `users` WHERE name LIKE "name%")
// SELECT
//   lk.article_id,
//   COUNT(lk.id) AS cnt
// FROM
//   likes lk
// GROUP BY
//   lk.article_id
// ;
