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
		articleGroup.GET("/:id", articleApi.GetById)
	}

	return router
}
