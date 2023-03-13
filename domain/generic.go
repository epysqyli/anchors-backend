package domain

import "gorm.io/gorm"

type Generic struct {
	gorm.Model
	Url   string `json:"url" gorm:"unique"`
	Ideas []Idea `gorm:"many2many:generics_ideas"`
}

type GenericsIdeas struct {
	IdeaID    uint
	GenericID uint
}
