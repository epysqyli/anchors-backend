package domain

import (
	"time"
)

type BlogsIdeas struct {
	IdeaID    uint `gorm:"primaryKey"`
	BlogID    uint `gorm:"primaryKey"`
	CreatedAt time.Time
}

func (BlogsIdeas) TableName() string {
	return "blogs_ideas"
}
