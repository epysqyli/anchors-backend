package domain

import "time"

type Author struct {
	ID             uint   `gorm:"primarykey"`
	OpenLibraryKey string `json:"open_library_key" gorm:"size:20;unique"`
	FullName       string `json:"full_name" gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
