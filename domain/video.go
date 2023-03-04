package domain

import (
	"strings"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Url            string `json:"url" gorm:"not null"`
	Identifier     string `gorm:"unique;not null"` // duplicate URL for non youtube videos for querying purposes
	YoutubeChannel string `json:"youtube_channel"` // retrieve via API if ab_channel or other ways not possible
	Ideas          []Idea `gorm:"many2many:ideas_videos;"`
	Timestamp      int16  `json:"timestamp" gorm:"-"`
}

// yt case: what are the potential url formats?
// youtube might have other short urls -> obtain full url before saving it
func (video *Video) AssignIdentifier() {
	if strings.Contains(video.Url, "youtube") {
		queryParams := strings.Split(video.Url, "?v=")
		restParams := strings.Split(queryParams[1], "&t")
		video.Identifier = restParams[0]
	} else if strings.Contains(video.Url, "youtu.be") {
		video.Identifier = strings.Split(video.Url, "youtu.be/")[1]
	} else {
		video.Identifier = video.Url
	}
}
