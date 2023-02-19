package domain

import (
	"time"
)

type IdeasVideos struct {
	IdeaID    uint `gorm:"primaryKey"`
	VideoID   uint `gorm:"primaryKey"`
	CreatedAt time.Time
}

func (IdeasVideos) TableName() string {
	return "ideas_videos"
}
