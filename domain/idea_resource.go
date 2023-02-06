package domain

import (
	"time"
)

type IdeaResource struct {
	IdeaId     uint `gorm:"primaryKey"`
	ResourceId uint `gorm:"primaryKey"`
	CreatedAt  time.Time
}

func (IdeaResource) TableName() string {
	return "ideas_resources"
}
