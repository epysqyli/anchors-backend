package repository

import (
	"context"

	"github.com/epysqyli/anchors-backend/domain"
	"gorm.io/gorm"
)

type IdeaRepository struct {
	database *gorm.DB
}

func NewIdeaRepository(db *gorm.DB) domain.IdeaRepository {
	return &IdeaRepository{
		database: db,
	}
}

func (ir *IdeaRepository) Create(c context.Context, idea *domain.Idea) error {
	ir.assignExistingIDs(idea)
	ir.assignResourceFields(idea)

	res := ir.database.Create(idea)
	if res.Error != nil {
		return res.Error
	}

	ir.assignRelationFields(idea)
	return nil
}

// the main feed: querying and pagination need to be set up
func (ir *IdeaRepository) FetchAll(c context.Context) ([]domain.Idea, error) {
	var ideas []domain.Idea
	res := ir.database.
		Preload("Blogs").
		Preload("Videos").
		Preload("Anchors").
		Preload("Books.Authors").
		Preload("Movies").
		Preload("Movies.Genres").
		Find(&ideas)

	return ideas, res.Error
}

func (ir *IdeaRepository) FetchByUserID(c context.Context, userID string) ([]domain.Idea, error) {
	var ideas []domain.Idea
	res := ir.database.Model(&domain.Idea{}).
		Preload("Blogs").
		Preload("Videos").
		Preload("Anchors").
		Preload("Books.Authors").
		Preload("Movies").
		Preload("Movies.Genres").
		Preload("Songs").
		Preload("Songs.MusicalAlbum").
		Preload("Songs.Artists").
		Find(&ideas, "user_id = ?", userID)

	return ideas, res.Error
}

// each 'preload' executes a query, optimize when and if necessary
func (ir *IdeaRepository) FetchByID(c context.Context, id string) (domain.Idea, error) {
	var idea domain.Idea
	res := ir.database.
		Preload("Blogs").
		Preload("Videos").
		Preload("Anchors").
		Preload("Books.Authors").
		Preload("Movies").
		Preload("Movies.Genres").
		First(&idea, id)

	return idea, res.Error
}

func (ir *IdeaRepository) DeleteByID(c context.Context, id string) error {
	var idea domain.Idea
	tx := ir.database.Delete(&idea, id)
	return tx.Error
}

func (ir *IdeaRepository) assignExistingIDs(idea *domain.Idea) {
	for iv, video := range idea.Videos {
		if video.ID == 0 {
			v := domain.Video{}
			ir.database.Where(&domain.Video{Identifier: video.RetrieveIdentifier()}).First(&v)

			if v.ID != 0 {
				videoPtr := &idea.Videos[iv]
				videoPtr.ID = v.ID
			}
		}
	}

	for ib, book := range idea.Books {
		if book.ID == 0 {
			b := domain.Book{}
			ir.database.Where(&domain.Book{OpenLibraryKey: book.OpenLibraryKey}).First(&b)

			if b.ID != 0 {
				bookPtr := &idea.Books[ib]
				bookPtr.ID = b.ID
			}

			if book.Authors == nil {
				continue
			}

			for ia, author := range book.Authors {
				a := domain.Author{}
				ir.database.Where(&domain.Author{OpenLibraryKey: author.OpenLibraryKey}).First(&a)

				if a.ID != 0 {
					authorPtr := &book.Authors[ia]
					authorPtr.ID = a.ID
				}
			}
		}
	}

	for im, movie := range idea.Movies {
		if movie.ID == 0 {
			m := domain.Movie{}
			ir.database.Where(&domain.Movie{Identifier: movie.Identifier}).First(&m)

			if m.ID != 0 {
				moviePtr := &idea.Movies[im]
				moviePtr.ID = m.ID
			}

			if movie.Genres == nil {
				continue
			}

			for ig, genre := range movie.Genres {
				g := domain.CinematicGenre{}
				ir.database.Where(&domain.CinematicGenre{Name: genre.Name}).First(&g)

				if g.ID != 0 {
					genrePtr := &movie.Genres[ig]
					genrePtr.ID = g.ID
				}
			}
		}
	}
}

// beforeCreate hook: assign unique identifiers and other computed fields based on the resource
func (ir *IdeaRepository) assignResourceFields(idea *domain.Idea) {
	// assign youtube channel as well - might just happen on the frontend
	for i, video := range idea.Videos {
		if video.Identifier == "" {
			videoPtr := &idea.Videos[i]
			videoPtr.AssignIdentifier()
		}
	}

	for i, song := range idea.Songs {
		songPtr := &idea.Songs[i]

		album := domain.MusicalAlbum{}
		ir.database.Where(&domain.MusicalAlbum{SpotifyID: song.MusicalAlbum.SpotifyID}).First(&album)
		if album.SpotifyID == "" {
			ir.database.Create(song.MusicalAlbum)
		}

		songPtr.MusicalAlbumSpotifyID = song.MusicalAlbum.SpotifyID
	}
}

// afterCreate DB hook?
func (ir *IdeaRepository) assignRelationFields(idea *domain.Idea) {
	for _, video := range idea.Videos {
		if video.Timestamp != 0 {
			ir.database.
				Model(domain.IdeasVideos{IdeaID: idea.ID, VideoID: video.ID}).
				Update("timestamp", video.Timestamp)
		}
	}

	for _, book := range idea.Books {
		if book.Chapter != "" {
			ir.database.
				Model(domain.BooksIdeas{IdeaID: idea.ID, BookID: book.ID}).
				Update("chapter", book.Chapter)
		}
	}

	for _, movie := range idea.Movies {
		if movie.Scene != "" {
			ir.database.
				Model(domain.IdeasMovies{IdeaID: idea.ID, MovieID: movie.ID}).
				Update("scene", movie.Scene)
		}
	}
}
