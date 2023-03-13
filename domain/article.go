package domain

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Url   string `json:"url" gorm:"unique"`
	Ideas []Idea `gorm:"many2many:articles_ideas"`
}

type ArticlesIdeas struct {
	IdeaID    uint
	ArticleID uint
}
