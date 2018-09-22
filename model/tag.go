package model

import (
	"errors"
	"time"
)

type Tag struct {
	Id        int `gorm:"auto_increment"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (Tag) SelectByArticleId(id int) ([]Tag, error) {
	tags := make([]Tag, 0)
	db.Raw("select maiev_tag.* from maiev_tag, maiev_article_tag "+
		"where maiev_article_tag.article_id = ? "+
		"and maiev_tag.id = maiev_article_tag.tag_id "+
		"and maiev_tag.deleted_at is null", id).
		Scan(&tags)
	return tags, db.Error
}

func (Tag) LinkTagsToArticle(articleId int, tags []Tag) error {
	if articleId == 0 {
		return errors.New("articleId is invalid")
	}
	tx := db.Begin()
	for _, tag := range tags {
		if tag.Name == "" {
			continue
		}
		var tagId int
		if err := tx.Table("maiev_tag").Select("id").Where("name = ?", tag.Name).Scan(&tagId).Error;
			err != nil {
			tx.Rollback()
			return err
		}
		if tagId == 0 {
			if err := tx.Create(tag).Error; err != nil {
				tx.Rollback()
				return err
			}
			tagId = tag.Id
		}
		if err := tx.Raw("insert into maiev_article_tag (tag_id, article_id) values (?, ?)", tagId, articleId).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}
