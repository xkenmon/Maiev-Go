package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/xkenmon/maiev/dto"
	"github.com/xkenmon/maiev/model"
	"strconv"
	"strings"
)

type ArticleApi struct{}

var articleModel = new(model.Article)

func (ArticleApi) GetById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(400, dto.ApiMessage{
			Code: 400,
			Msg:  "id转换错误，请确定id为整数",
		})
		return
	}
	var withContent bool
	withContentStr := c.Query("content")
	if withContentStr != "" {
		withContent, err = strconv.ParseBool(withContentStr)
		if err != nil {
			c.AbortWithStatusJSON(400, dto.ApiMessage{
				Code: 400,
				Msg:  "content转换错误，请确定content为bool",
			})
		}
	}
	article, err := articleModel.GetById(id, withContent)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(200, dto.ApiMessage{
				Code: 200,
				Msg:  err.Error(),
			})
		default:
			c.AbortWithStatusJSON(500, dto.ApiMessage{
				Code: 500,
				Msg:  "服务器内部错误:" + err.Error(),
			})
		}
	}
	c.JSON(200, article)
}

func (ArticleApi) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.AbortWithStatusJSON(400, dto.ApiMessage{
			Code: 400,
			Msg:  "param page except int:" + err.Error(),
		})
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.AbortWithStatusJSON(400, dto.ApiMessage{
			Code: 400,
			Msg:  "param limit except int:" + err.Error(),
		})
	}
	order := c.DefaultQuery("order", "id")
	sort := c.DefaultQuery("sort", "asc")

	sortValid := strings.EqualFold(sort, "asc") || strings.EqualFold(sort, "desc")

	if !sortValid {
		c.AbortWithStatusJSON(400, dto.ApiMessage{
			Code: 400,
			Msg:  "The sort field should be asc or desc",
		})
	}

	articles, err := articleModel.List(page, limit, order, sort)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(200, dto.ApiMessage{
				Code: 200,
				Msg:  err.Error(),
			})
		default:
			c.AbortWithStatusJSON(500, dto.ApiMessage{
				Code: 500,
				Msg:  "内部错误:" + err.Error(),
			})
		}
	}
	c.AbortWithStatusJSON(200, articles)
}
