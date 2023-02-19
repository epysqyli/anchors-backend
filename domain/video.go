package domain

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Url            string `json:"url" form:"url"`
	YoutubeChannel string `json:"youtube_channel"`
}
