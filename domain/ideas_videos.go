package domain

type IdeasVideos struct {
	IdeaID    uint  `gorm:"primaryKey"`
	VideoID   uint  `gorm:"primaryKey"`
	Timestamp int16 `json:"timestamp"`
}

func (IdeasVideos) TableName() string {
	return "ideas_videos"
}
