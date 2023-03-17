package domain

type Tag struct {
	ID   uint   `gorm:"primarykey"`
	Name string `json:"chapter" gorm:"not null;size:100"`
}

type TagRepository interface {
	Create(tag *Tag) error
	FetchAll() []Tag
	FetchById(ID string) Tag
	FetchByName(name string) Tag
	DeleteByID(ID string) error
}
