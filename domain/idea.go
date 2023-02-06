package domain

import (
	"context"

	"gorm.io/gorm"
)

type Idea struct {
	gorm.Model
	UserId    uint       `json:"user_id" form:"user_id"`
	Content   string     `json:"content" form:"content"`
	Resources []Resource `gorm:"many2many:ideas_resources;" json:"resources" form:"resources"`
}

type IdeaRepository interface {
	Create(c context.Context, idea *Idea) error
	FetchByUserID(c context.Context, userId string) ([]Idea, error)
}

type IdeaUsecase interface {
	Create(c context.Context, idea *Idea) error
	FetchByUserID(c context.Context, userId string) ([]Idea, error)
}
