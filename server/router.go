package server

import (
	"github.com/gin-gonic/gin"
	"github.com/xkenmon/maiev/api"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	articleGroup := router.Group("/article")
	{
		articleApi := new(api.ArticleApi)
		articleGroup.GET("", articleApi.List)
		articleGroup.POST("", articleApi.Insert)
		articleGroup.GET("/tag-name/:tag", articleApi.SelectByTagName)
		articleGroup.GET("/:id", articleApi.GetById)
	}
	commentGroup := router.Group("/comment")
	{
		commentApi := new(api.CommentApi)
		commentGroup.GET("article/:articleId", commentApi.ListByArticleId)
		commentGroup.POST("article", commentApi.AddArticleComment)
		commentGroup.GET("parent/:commentId", commentApi.ListByParentId)
		commentGroup.POST("parent", commentApi.AddParentComment)
	}

	return router
}
