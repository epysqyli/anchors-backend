package domain

import "gorm.io/gorm"

type Wiki struct {
	gorm.Model
	Url   string `json:"url" gorm:"unique"`
	Ideas []Idea `gorm:"many2many:ideas_wikis"`
}

type IdeasWikis struct {
	IdeaID uint
	WikiID uint
}
