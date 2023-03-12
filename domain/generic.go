package domain

import "gorm.io/gorm"

type Generic struct {
	gorm.Model
	Url string `json:"url" gorm:"unique"`
}

type GenericsIdeas struct {
	IdeaID    uint
	GenericID uint
}
