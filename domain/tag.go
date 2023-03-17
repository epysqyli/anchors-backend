package domain

type Tag struct {
	ID   uint   `gorm:"primarykey"`
	Name string `json:"chapter" gorm:"not null;size:100"`
}
