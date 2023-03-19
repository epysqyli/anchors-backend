package domain

type Tag struct {
	ID    uint   `gorm:"primarykey"`
	Name  string `json:"chapter" gorm:"not null;size:100;unique"`
	Ideas []Idea `json:"ideas" gorm:"many2many:ideas_tags"`
}

type TagIdeasRequest struct {
	ID  uint   `json:"id"`
	And []uint `json:"and"`
	Or  []uint `json:"or"`
	Not []uint `json:"not"`
}

type TagRepository interface {
	Create(tag *Tag) error
	FetchAll() []Tag
	FetchById(ID string) Tag
	FetchByName(name string) Tag
	DeleteByID(ID string) error
}
