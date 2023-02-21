package domain

import (
	"context"

	"gorm.io/gorm"
)

type Idea struct {
	gorm.Model
	UserID  uint    `json:"user_id"`
	Content string  `json:"content"`
	Videos  []Video `gorm:"many2many:ideas_videos;" json:"videos"`
	Blogs   []Blog  `gorm:"many2many:blogs_ideas;" json:"blogs"`
}

type IdeaRepository interface {
	Create(c context.Context, idea *Idea) error
	FetchByUserID(c context.Context, userId string) ([]Idea, error)
	FetchByID(c context.Context, id string) (Idea, error)
	FetchAll(c context.Context) ([]Idea, error)
	DeleteByID(c context.Context, id string) error
}

func (idea Idea) HasNoResources() bool {
	if idea.Blogs == nil && idea.Videos == nil {
		return true
	}

	return false
}
