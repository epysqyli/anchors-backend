package domain

import "time"

type BooksIdeas struct {
	IdeaID    uint   `gorm:"primaryKey"`
	BookID    uint   `gorm:"primaryKey"`
	Chapter   string `json:"chapter" gorm:"size:256"`
	CreatedAt time.Time
}

func (BooksIdeas) TableName() string {
	return "books_ideas"
}
