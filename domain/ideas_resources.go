package domain

import (
	"time"
)

type IdeasResources struct {
	IdeaId     uint `gorm:"primaryKey"`
	ResourceId uint `gorm:"primaryKey"`
	CreatedAt  time.Time
}

func (IdeasResources) TableName() string {
	return "ideas_resources"
}
