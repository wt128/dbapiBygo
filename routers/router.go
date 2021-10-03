package routers

import (
	"ggapi/controller/like"
	"ggapi/controller/post"
	"ggapi/controller/session"
	"ggapi/controller/user"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func Init() {
	r := router()
	//r.RunTLS(":8080", "./server.pem", "./server.key")
	r.Run(":8080")
}

func router() *gin.Engine {
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	//r.Use(static.Serve("/",static.LocalFile("./view",true)))
	r.Static("/view","./view")
	config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://localhost:3000"}
    r.Use(cors.New(config))
	
	u := r.Group("/users")
	{
		ctrl := user.Controller{}
        u.POST("", ctrl.Create)
		u.POST("/setimg",ctrl.Setimg)
		u.GET("/getimg/:id",ctrl.Getimg)
        u.POST("/update", ctrl.Update)
        u.POST("/destory/:id", ctrl.Delete)
		u.GET("/:id",ctrl.Index)
		
	}

	s := r.Group("/session")
	{
		ctrl := session.Controller{}
		s.POST("",ctrl.SetId)
		s.POST("/create",ctrl.Create)
		s.POST("/email",ctrl.GetEmail)
		// s.POST("/destroy",) logout処理
	}

	p := r.Group("/post")
	{
		ctrl := post.Controller{}
		p.POST("/new",ctrl.Post) 
		p.POST("/destroy",ctrl.Delete)
		p.GET("/show",ctrl.Show) //記事一覧
		p.GET("/total",ctrl.PostTotal)
		p.GET("/:id",ctrl.GetOne)
		p.POST("/setimg/:pid",ctrl.Setimg)
		p.GET("/getimg/:pid",ctrl.Getimg)
		p.POST("/update",ctrl.Update)
	}

	l := r.Group("/like")
	{
		ctrl := like.Controller{}
		l.GET("",ctrl.Show)
		l.GET("/:id",ctrl.Auser)
		l.GET("/check/:id",ctrl.Check)
		l.POST("/create",ctrl.Create)
		l.POST("/destroy",ctrl.Destroy)
	}
	return r
}
