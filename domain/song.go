package domain

import "time"

type Song struct {
	SpotifyID             string          `json:"spotify_id" gorm:"primarykey;size:100"`
	Name                  string          `json:"name" gorm:"size:512"`
	SpotifyUrl            string          `json:"spotify_url" gorm:"size:512"`
	PreviewUrl            string          `json:"preview_url" gorm:"size:512"`
	Artists               []MusicalArtist `json:"artists" gorm:"many2many:musical_artists_songs"`
	Album                 MusicalAlbum    `json:"album" gorm:"-"`
	MusicalAlbumSpotifyID string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type MusicalAlbum struct {
	SpotifyID   string `json:"spotify_id" gorm:"primarykey;size:100"`
	SpotifyUrl  string `json:"spotify_url" gorm:"size:512"`
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date" gorm:"size:10"`
	CoverUrl    string `json:"cover_url" gorm:"size:512"`
	Songs       []Song
}

type MusicalArtist struct {
	SpotifyID  string `json:"spotify_id" gorm:"primarykey;size:100"`
	SpotifyUrl string `json:"spotify_url"`
	Name       string `json:"name"`
}

type MusicalArtistsSongs struct {
	SongSpotifyID          string
	MusicalArtistSpotifyID string
}

type IdeasSongs struct {
	SongSpotifyID string
	IdeaID        uint
}
