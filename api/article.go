package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/xkenmon/maiev/model"
	"strconv"
	"strings"
)

type ArticleApi struct{}

var articleModel = new(model.Article)

func (ArticleApi) GetById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if checkErrAndWrite(c, err, 400, fmt.Sprintf("param id should be int,but %s find.", c.Param("id"))) {
		return
	}
	var withContent bool
	withContentStr := c.Query("content")
	if withContentStr != "" {
		withContent, err = strconv.ParseBool(withContentStr)
		if checkErrAndWrite(c, err, 400, "param content should be bool,but "+withContentStr+" find") {
			return
		}
	}
	article, err := articleModel.GetById(id, withContent)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			writeMsg(c, 200, "record not found")
			return
		default:
			writeMsg(c, 500, "internal error.")
			return
		}
	}
	c.JSON(200, article)
}

func (ArticleApi) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if checkErrAndWrite(c, err, 400, "param page should be int, but "+c.Query("page")+" found.") {
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if checkErrAndWrite(c, err, 400, "param limit should be int, but "+c.Query("limit")+" found.") {
		return
	}
	order := c.DefaultQuery("order", "id")
	sort := c.DefaultQuery("sort", "asc")

	sortValid := strings.EqualFold(sort, "asc") || strings.EqualFold(sort, "desc")

	if !sortValid {
		writeMsg(c, 400, "The sort field should be asc or desc")
		return
	}

	articles, err := articleModel.List(page, limit, order, sort)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			writeMsg(c, 200, "record not found")
			return
		default:
			writeMsg(c, 500, "internal error.")
			return
		}
	}
	c.JSON(200, articles)
}
