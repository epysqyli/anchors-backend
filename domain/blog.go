package domain

import "gorm.io/gorm"

type Blog struct {
	gorm.Model
	Url      string `json:"url" gorm:"not null;unique"`
	Category string `json:"category"`
	Ideas    []Idea `gorm:"many2many:blogs_ideas;"`
}
