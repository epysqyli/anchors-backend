package domain

import "gorm.io/gorm"

type Blog struct {
	gorm.Model
	Url      string `json:"url" gorm:"unique_index:blog_url"`
	Category string `json:"category"`
}
