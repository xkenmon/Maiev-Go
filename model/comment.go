package model

import "time"

type CommentArticle struct {
	Id          int `gorm:"auto_increment"`
	ArticleId   int
	IsAnonymous bool
	UserId      int
	UserName    string `gorm:"-"`
	UserAvt     string `gorm:"-"`
	Content     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type CommentSelf struct {
	Id          int `gorm:"auto_increment"`
	ArticleId   int
	ParentId    int
	IsAnonymous bool
	UserId      int
	UserName    string `gorm:"-"`
	UserAvt     string `gorm:"-"`
	Content     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// ListByArticleId select article comment by articleId
func (CommentArticle) ListByArticleId(articleId, page, limit int, order, sort string) []CommentArticle {
	comments := make([]CommentArticle, 0)
	db.Raw("select c.*, u.name as user_name, u.avatar as user_avt "+
		"from maiev_comment_article c left join maiev_user u on c.user_id = u.id "+
		"where c.article_id = ? and c.deleted_at is null ", articleId).
		Limit(limit).
		Offset((page - 1) * limit).
		Order(order + " " + sort).
		Scan(&comments)
	for _, c := range comments {
		if c.IsAnonymous {
			c.UserId = 0
			c.UserAvt = ""
			c.UserName = "匿名用户"
		}
	}
	return comments
}

// ListByCommentId select comment's comment by parent comment id
func (CommentSelf) ListByCommentId(commentId, page, limit int, order, sort string) []CommentSelf {
	comments := make([]CommentSelf, 0)
	db.Raw("select c.*, u.name as user_name, u.avatar as user_avt "+
		"from maiev_comment_self c left join maiev_user u on c.user_id = u.id "+
		"where c.parent_id = ? and c.deleted_at is null ", commentId).
		Order(order + " " + sort).
		Limit(limit).
		Offset((page - 1) * limit).
		Scan(&comments)
	for _, c := range comments {
		if c.IsAnonymous {
			c.UserId = 0
			c.UserName = "匿名用户"
			c.UserAvt = ""
		}
	}
	return comments
}

func (CommentArticle) InsertComment(comment *CommentArticle) *CommentArticle {
	db.Create(comment)
	return comment
}

func (CommentSelf) InsertComment(comment *CommentSelf) *CommentSelf {
	db.Create(comment)
	return comment
}
