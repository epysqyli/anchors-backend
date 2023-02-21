package domain

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Url            string `json:"url" gorm:"not null;unique"`
	YoutubeChannel string `json:"youtube_channel"`
}
