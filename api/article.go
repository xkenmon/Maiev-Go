package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/xkenmon/maiev/dto"
	"github.com/xkenmon/maiev/model"
	"strconv"
)

type ArticleApi struct{}

const ArticleOrder string = "created_at|updated_at|like_count|read_count|title|id"

var articleModel = new(model.Article)

func (ArticleApi) GetById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeMsg(c, 400, fmt.Sprintf("param id should be int,but %s find.", c.Param("id")), nil)
		return
	}

	var withContent bool

	if withContentStr := c.Query("content"); withContentStr != "" {
		withContent, err = strconv.ParseBool(withContentStr)
		if err != nil {
			writeMsg(c, 400, "param content should be bool,but "+withContentStr+" find", nil)
			return
		}
	}

	if article, err := articleModel.GetById(id, withContent); err != nil {
		if err == gorm.ErrRecordNotFound {
			writeMsg(c, 404, "record not found", nil)
			return
		}
		writeMsg(c, 500, err.Error(), nil)
		return
	} else {
		writeMsg(c, 200, "ok", article)
	}
}

func (ArticleApi) SelectByTagName(c *gin.Context) {
	tagName := c.Param("tag")
	if tagName == "" {
		writeMsg(c, 400, "未指定tag", nil)
		return
	}
	pageReq := dto.PageRequest{Page: 1, Limit: 10, Order: "created_at", Sort: "asc"}

	if err := c.Bind(&pageReq); err != nil {
		writeMsg(c, 400, err.Error(), nil)
		return
	}

	if !isOrderValid(pageReq.Order, ArticleOrder) || !isSortValid(pageReq.Sort) {
		writeMsg(c, 400, "invalid order or sort field", nil)
		return
	}

	if articles, err := articleModel.SelectByTagName(tagName); err != nil {
		writeMsg(c, 500, err.Error(), nil)
	} else {
		writeMsg(c, 200, "ok", articles)
	}
}

func (ArticleApi) List(c *gin.Context) {
	pageReq := dto.PageRequest{Page: 1, Limit: 10, Order: "created_at", Sort: "asc"}

	if err := c.Bind(&pageReq); err != nil {
		writeMsg(c, 400, err.Error(), nil)
		return
	}

	if !isOrderValid(pageReq.Order, ArticleOrder) || !isSortValid(pageReq.Sort) {
		writeMsg(c, 400, "invalid order or sort field", nil)
		return
	}

	if articles, err := articleModel.List(pageReq.Page, pageReq.Limit, pageReq.Order, pageReq.Sort); err != nil {
		writeMsg(c, 500, "internal server error:"+err.Error(), nil)
		return
	} else {
		writeMsg(c, 200, "ok", articles)
	}
}

func (ArticleApi) Insert(c *gin.Context) {
	article := &model.Article{}
	if article.UserId == 0 || article.Title == "" || article.Content == "" {
		writeMsg(c, 400, "字段错误", nil)
		return
	}
	if err := articleModel.Insert(article); err != nil {
		writeMsg(c, 200, err.Error(), nil)
		return
	}
	writeMsg(c, 200, "ojbk", article)
}
