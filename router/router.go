package router

import (
	"github.com/gin-gonic/gin"
	"myDouYIn/controller"
)

func InitRouter() {
	r := gin.Default()

	v1 := r.Group("/douyin")

	v1.GET("/user/", controller.Action)
	v1.POST("/user/register/", controller.Register)
	v1.POST("/user/login/", controller.Login)

	v1.POST("/publish/action/", controller.PublishLogin)
	v1.GET("/publish/list/", controller.PublishList)

	v1.GET("/feed/", controller.FeedAction)

	v1.GET("/video/", controller.FeedVideo)

	v1.POST("favorite/action/", controller.FavoriteAction)
	v1.GET("favorite/list/", controller.FavoriteList)

	v1.POST("comment/action/", controller.CommentAction)
	v1.GET("comment/list/", controller.CommentList)

	v1.POST("relation/action/", controller.RelationAction)
	v1.GET("relation/follow/list", controller.RelationFollowList)
	v1.GET("relation/follower/list", controller.RelationFollowerList)
	r.Run(":2222")
}
