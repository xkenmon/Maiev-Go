package model

import (
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type Article struct {
	Id          int    `gorm:"auto_increment"`
	Title       string `gorm:"size:100"`
	UserId      uint
	IsAnonymous bool
	LikeCount   int
	ReadCount   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ContentId   int
	Content     string `gorm:"-"`
	Summary     string `gorm:"-"`
	Tags        []Tag  `gorm:"-"`
	DeletedAt   *time.Time
}

type ArticleContent struct {
	Id        int    `gorm:"auto_increment"`
	Content   string `gorm:"type:longtext"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

const ArticleSummaryLen = 100

var tagModel = new(Tag)

func (Article) GetById(id int, withContent bool) (*Article, error) {
	var article Article

	if err := db.First(&article, id).Error; err != nil {
		return nil, err
	}

	tags, err := tagModel.SelectByArticleId(id)
	if err != nil {
		log.Print(err)
	}
	article.Tags = tags
	if article.ContentId != 0 {
		if withContent {
			var articleContent ArticleContent
			db.First(&articleContent, article.ContentId)
			article.Content = articleContent.Content
			if len(article.Content) <= ArticleSummaryLen {
				article.Summary = article.Content
			} else {
				article.Summary = article.Content[0:ArticleSummaryLen]
			}
		} else {
			var articleSummary ArticleContent
			db.Select("SUBSTRING(content,1,100) as content").First(&articleSummary, article.ContentId)
			article.Summary = articleSummary.Content
		}
	}
	return &article, nil
}

func (Article) SelectByTagName(tagName string) ([]Article, error) {
	articles := make([]Article, 0)
	db.Raw("select maiev_article.* "+
		"from maiev_article, maiev_article_tag, maiev_tag "+
		"where maiev_article_tag.tag_id = maiev_tag.id "+
		"and maiev_tag.name = ? "+
		"and maiev_article.id = maiev_article_tag.article_id "+
		"and maiev_tag.deleted_at is null "+
		"and maiev_article.deleted_at is null", tagName).Scan(&articles)
	return articles, db.Error
}

func (Article) SelectByTagId(id int) []Article {
	articles := make([]Article, 0)
	db.Raw("select maiev_article.* "+
		"from maiev_article, maiev_article_tag, maiev_article_tag "+
		"where maiev_article_tag.tag_id = ? "+
		"and maiev_article.id = maiev_article_tag.article_id "+
		"and maiev_tag.deleted_at is null "+
		"and maiev_article.deleted_at is null", id).Scan(&articles)
	return articles
}

func (Article) List(page int, limit int, order string, sort string) (articles []*Article, err error) {
	articles = make([]*Article, 0)
	db.Limit(limit).
		Offset((page - 1) * limit).
		Order(order + " " + sort).
		Find(&articles)
	if len(articles) == 0 {
		return articles, gorm.ErrRecordNotFound
	}
	return articles, db.Error
}

func (Article) Insert(article *Article) error {
	tx := db.Begin()
	// 插入文章
	if err := tx.Create(article).Error; err != nil {
		tx.Rollback()
		return err
	}
	//插入文章内容
	content := &ArticleContent{Content: article.Content}
	tx.Create(content)
	tx.Raw("update maiev_article "+
		"set content_id = ? "+
		"where maiev_article.id = ?", content.Id, article.Id)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	//文章标签处理
	if err := tagModel.LinkTagsToArticle(article.Id, article.Tags); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
