package api

import (
	"github.com/gin-gonic/gin"
	"github.com/xkenmon/maiev/dto"
	"github.com/xkenmon/maiev/model"
	"strconv"
)

const CommentOrder string = "id|created_at|updated_at"

type CommentApi struct{}

var commentArticleModel = new(model.CommentArticle)

var commentSelfModel = new(model.CommentSelf)

func (CommentApi) ListByArticleId(c *gin.Context) {
	articleId, err := strconv.Atoi(c.Param("articleId"))
	if err != nil {
		writeMsg(c, 400, "param articleId should be int", nil)
		return
	}

	pageReq := dto.PageRequest{Page: 1, Limit: 10, Order: "created_at", Sort: "asc"}

	if err := c.Bind(&pageReq); err != nil {
		writeMsg(c, 400, err.Error(), nil)
		return
	}

	if !isOrderValid(pageReq.Order, CommentOrder) || !isSortValid(pageReq.Sort) {
		writeMsg(c, 400, "invalid order or sort field", nil)
		return
	}

	result := commentArticleModel.ListByArticleId(articleId, pageReq.Page, pageReq.Limit, pageReq.Order, pageReq.Sort)
	writeMsg(c, 200, "ojbk", result)
}

func (CommentApi) ListByParentId(c *gin.Context) {
	commentId, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		writeMsg(c, 400, err.Error(), nil)
		return
	}

	pageReq := dto.PageRequest{Page: 1, Limit: 10, Order: "created_at", Sort: "asc"}

	if err := c.Bind(&pageReq); err != nil {
		writeMsg(c, 400, err.Error(), nil)
		return
	}

	if !isOrderValid(pageReq.Order, CommentOrder) || !isSortValid(pageReq.Sort) {
		writeMsg(c, 400, "invalid order or sort field", nil)
		return
	}

	result := commentSelfModel.ListByCommentId(commentId, pageReq.Page, pageReq.Limit, pageReq.Order, pageReq.Sort)

	writeMsg(c, 200, "okay", result)
}

func (CommentApi) AddArticleComment(c *gin.Context) {
	comment := &model.CommentArticle{}
	if err := c.BindJSON(comment); err != nil {
		writeMsg(c, 400, err.Error(), nil)
		return
	}
	if comment.UserId == 0 {
		writeMsg(c, 400, "userId is invalid", nil)
		return
	}
	if comment.Content == "" {
		writeMsg(c, 400, "content is empty", nil)
		return
	}

	writeMsg(c, 200, "ok", commentArticleModel.InsertComment(comment))
}

func (CommentApi) AddParentComment(c *gin.Context) {
	comment := &model.CommentSelf{}
	if err := c.BindJSON(comment); err != nil {
		writeMsg(c, 400, err.Error(), nil)
		return
	}
	if comment.UserId == 0 || comment.ParentId == 0 {
		writeMsg(c, 400, "userId or parentId is invalid", nil)
		return
	}
	if comment.Content == "" {
		writeMsg(c, 400, "content is empty", nil)
		return
	}

	writeMsg(c, 200, "ok", commentSelfModel.InsertComment(comment))
}
