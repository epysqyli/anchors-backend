package domain

import "gorm.io/gorm"

// have also ideas_movies struct on this same file for clarity?

type Movie struct {
	gorm.Model
	Identifier       uint32           `json:"identifier" gorm:"not null;unique"` // id from TDBM frontend resp
	Title            string           `json:"title" gorm:"not null"`
	OriginalTitle    string           `json:"original_title" gorm:"not null"`
	PosterPath       string           `json:"poster_path" gorm:"size:100"`
	ReleaseDate      string           `json:"release_date" gorm:"size:10"`
	Runtime          uint16           `json:"runtime"`
	OriginalLanguage string           `json:"original_language" gorm:"not null;size:30"`
	Genres           []CinematicGenre `json:"genres" gorm:"many2many:cinematic_genres_movies"`
}

type CinematicGenre struct {
	ID   uint   `gorm:"primarykey"`
	Name string `json:"name" gorm:"not null;size:255;unique"`
}
