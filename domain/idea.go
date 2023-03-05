package domain

import (
	"context"

	"gorm.io/gorm"
)

// which associations should be array of pointers?
type Idea struct {
	gorm.Model
	UserID  uint    `json:"user_id"`
	Content string  `json:"content"`
	Videos  []Video `gorm:"many2many:ideas_videos;" json:"videos"`
	Blogs   []Blog  `gorm:"many2many:blogs_ideas;" json:"blogs"`
	Books   []Book  `gorm:"many2many:books_ideas" json:"books"`
	Movies  []Movie `gorm:"many2many:ideas_movies" json:"movies"`
	Anchors []*Idea `gorm:"many2many:anchors_ideas;" json:"anchors"`
}

/**
 * add optional arg: withAssociations?
 * are there gonna be cases where we want only to fetch the bare idea model?
 */
type IdeaRepository interface {
	Create(c context.Context, idea *Idea) error
	FetchByUserID(c context.Context, userID string) ([]Idea, error)
	FetchByID(c context.Context, id string) (Idea, error)
	FetchAll(c context.Context) ([]Idea, error)
	DeleteByID(c context.Context, id string) error
}

func (idea Idea) HasNoResources() bool {
	if idea.Blogs == nil &&
		idea.Videos == nil &&
		idea.Books == nil &&
		idea.Movies == nil &&
		idea.Anchors == nil {
		return true
	}

	return false
}
