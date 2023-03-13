package domain

import (
	"context"

	"gorm.io/gorm"
)

type Idea struct {
	gorm.Model
	UserID   uint      `json:"user_id"`
	Content  string    `json:"content"`
	Videos   []Video   `gorm:"many2many:ideas_videos;" json:"videos"`
	Blogs    []Blog    `gorm:"many2many:blogs_ideas;" json:"blogs"`
	Books    []Book    `gorm:"many2many:books_ideas" json:"books"`
	Movies   []Movie   `gorm:"many2many:ideas_movies" json:"movies"`
	Songs    []Song    `gorm:"many2many:ideas_songs" json:"songs"`
	Wikis    []Wiki    `gorm:"many2many:ideas_wikis" json:"wikis"`
	Generics []Generic `gorm:"many2many:generics_ideas" json:"generics"`
	Articles []Article `gorm:"many2many:articles_ideas" json:"articles"`
	Anchors  []Idea    `gorm:"many2many:anchors_ideas;" json:"anchors"` // array of pointers?
}

/**
 * add optional arg: withAssociations?
 * are there gonna be cases where we want only to fetch the bare idea model?
 */
type IdeaRepository interface {
	Create(c context.Context, idea *Idea) error
	FetchByUserID(c context.Context, userID string) ([]Idea, error)
	FetchByID(c context.Context, id string) (Idea, error)
	FetchByResourceID(c context.Context, resType string, resID string) []Idea
	FetchAll(c context.Context) ([]Idea, error)
	DeleteByID(c context.Context, id string) error
}

func (idea Idea) HasNoResources() bool {
	if idea.Blogs == nil &&
		idea.Videos == nil &&
		idea.Books == nil &&
		idea.Movies == nil &&
		idea.Songs == nil &&
		idea.Wikis == nil &&
		idea.Generics == nil &&
		idea.Articles == nil &&
		idea.Anchors == nil {
		return true
	}

	return false
}
