package domain

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Url string `json:"url" gorm:"unique"`
}

type ArticlesIdeas struct {
	IdeaID    uint
	ArticleID uint
}
