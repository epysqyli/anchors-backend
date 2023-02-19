package domain

import "gorm.io/gorm"

type VideoResource struct {
	gorm.Model
	ResourceID     uint   `json:"resource_id"`
	YoutubeChannel string `json:"youtube_channel"` // to be analyzed
}

func (VideoResource) TableName() string {
	return "videos"
}
