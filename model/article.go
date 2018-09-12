package model

import (
	"github.com/jinzhu/gorm"
	"github.com/xkenmon/maiev/database"
	"time"
)

type Article struct {
	Id          int    `gorm:"auto_increment"`
	Title       string `gorm:"size:100"`
	UserId      uint
	IsAnonymous bool
	LikeCount   int
	ReadCount   int
	GmtCreate   time.Time
	GmtUpdate   time.Time
	ContentId   int
	Content     string
	IsDelete    bool
}

type ArticleContent struct {
	Id      int    `gorm:"auto_increment"`
	Content string `gorm:"type:longtext"`
}

var db = database.GetDB()

func (Article) GetById(id int, withContent bool) (*Article, error) {
	var article Article
	db.Where("id = ? and is_delete = ?", id, false).First(&article)
	if article.Id == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if withContent && article.ContentId != 0 {
		var articleContent ArticleContent
		db.First(&articleContent, article.ContentId)
		article.Content = articleContent.Content
	}
	return &article, nil
}

func (Article) List(page int, limit int, order string, sort string) (articles []*Article, err error) {
	articles = make([]*Article, 0)
	db.Where("is_delete = false").
		Limit(limit).
		Offset((page - 1) * limit).
		Order(order + " " + sort).
		Find(&articles)
	if len(articles) == 0 {
		return articles, gorm.ErrRecordNotFound
	}
	return articles, nil
}
