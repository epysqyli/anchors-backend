package domain

import "context"

type Tag struct {
	ID   uint   `gorm:"primarykey"`
	Name string `json:"chapter" gorm:"not null;size:100"`
}

type TagRepository interface {
	Create(c context.Context, tag *Tag) (Tag, error)
	FetchById(c context.Context, ID string) (Tag, error)
	FetchByName(c context.Context, name string) (Tag, error)
	Delete(c context.Context, ID string) (Tag error)
}
