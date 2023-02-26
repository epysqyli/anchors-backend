package domain

import "time"

type Author struct {
	ID        uint   `gorm:"primarykey"`
	Key       string `json:"author_key" gorm:"size:20;unique"`
	FullName  string `json:"full_name" gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
