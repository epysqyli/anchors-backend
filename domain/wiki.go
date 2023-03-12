package domain

import "gorm.io/gorm"

type Wiki struct {
	gorm.Model
	Url string `json:"url" gorm:"unique"`
}

type IdeasWikis struct {
	IdeaID uint
	WikiID uint
}
