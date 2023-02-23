package domain

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Url            string `json:"url" gorm:"not null;unique"`
	YoutubeChannel string `json:"youtube_channel"`
	Ideas          []Idea `gorm:"many2many:ideas_videos;"`
	Timestamp      int16  `json:"timestamp" gorm:"-"`
}
