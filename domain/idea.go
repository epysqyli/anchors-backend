package domain

import (
	"context"

	"gorm.io/gorm"
)

type Idea struct {
	gorm.Model
	UserID  uint    `json:"user_id" form:"user_id"`
	Content string  `json:"content" form:"content"`
	Videos  []Video `gorm:"many2many:ideas_videos;" json:"videos" form:"videos"`
	Blogs   []Blog  `gorm:"many2many:blogs_ideas;" json:"blogs" form:"blogs"`
}

type IdeaRepository interface {
	Create(c context.Context, idea *Idea) error
	FetchByUserID(c context.Context, userId string) ([]Idea, error)
	FetchByID(c context.Context, id string) (Idea, error)
	FetchAll(c context.Context) ([]Idea, error)
	DeleteByID(c context.Context, id string) error
}
