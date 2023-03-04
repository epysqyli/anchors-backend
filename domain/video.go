package domain

import (
	"strings"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Url            string `json:"url" gorm:"not null"`
	Identifier     string `json:"identifier" gorm:"unique;not null"` // duplicate URL for non youtube videos for querying purposes
	YoutubeChannel string `json:"youtube_channel"`                   // gorm:"-"` retrieve via API if ab_channel or alternative missing
	Ideas          []Idea `gorm:"many2many:ideas_videos;"`
	Timestamp      int16  `json:"timestamp" gorm:"-"`
}

func (video Video) RetrieveIdentifier() string {
	if strings.Contains(video.Url, "youtube") {
		queryParams := strings.Split(video.Url, "?v=")
		restParams := strings.Split(queryParams[1], "&t")
		return restParams[0]
	}

	if strings.Contains(video.Url, "youtu.be") {
		return strings.Split(video.Url, "youtu.be/")[1]
	}

	return video.Url
}

// yt case: what are the potential url formats?
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
