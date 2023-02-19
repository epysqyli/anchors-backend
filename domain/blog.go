package domain

import "gorm.io/gorm"

type Blog struct {
	gorm.Model
	Url      string `json:"url"`
	Category string `json:"category"`
}
