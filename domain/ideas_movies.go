package domain

type IdeasMovies struct {
	IdeaID  uint   `gorm:"primaryKey"`
	MovieID uint   `gorm:"primaryKey"`
	Scene   string `json:"scene" gorm:"size:256"`
}

func (IdeasMovies) TableName() string {
	return "ideas_movies"
}
