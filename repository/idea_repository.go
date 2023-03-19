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

func (ir *IdeaRepository) FetchAll(c context.Context) ([]domain.Idea, error) {
	var ideas []domain.Idea
	res := ir.database.
		Preload("Blogs").
		Preload("Videos").
		Preload("Anchors").
		Preload("Books.Authors").
		Preload("Movies").
		Preload("Movies.Genres").
		Preload("Songs").
		Preload("Songs.MusicalAlbum").
		Preload("Songs.Artists").
		Preload("Wikis").
		Preload("Generics").
		Preload("Articles").
		Find(&ideas)

	return ideas, res.Error
}

func (ir *IdeaRepository) FetchByUserID(c context.Context, userID string) ([]domain.Idea, error) {
	var ideas []domain.Idea
	res := ir.database.
		Preload("Blogs").
		Preload("Videos").
		Preload("Anchors").
		Preload("Books.Authors").
		Preload("Movies").
		Preload("Movies.Genres").
		Preload("Songs").
		Preload("Songs.MusicalAlbum").
		Preload("Songs.Artists").
		Preload("Wikis").
		Preload("Generics").
		Preload("Articles").
		Find(&ideas, "user_id = ?", userID)

	return ideas, res.Error
}

func (ir *IdeaRepository) FetchByID(c context.Context, id string) (domain.Idea, error) {
	var idea domain.Idea
	res := ir.database.
		Preload("Blogs").
		Preload("Videos").
		Preload("Anchors").
		Preload("Books.Authors").
		Preload("Movies").
		Preload("Movies.Genres").
		Preload("Songs").
		Preload("Songs.MusicalAlbum").
		Preload("Songs.Artists").
		Preload("Wikis").
		Preload("Generics").
		Preload("Articles").
		First(&idea, id)

	return idea, res.Error
}

func (ir *IdeaRepository) FetchByTags(tagReq domain.TagQuery) (domain.Tag, error) {
	tag := domain.Tag{}

	ir.database.
		Preload("Ideas.Blogs").
		Preload("Ideas.Videos").
		Preload("Ideas.Anchors").
		Preload("Ideas.Books").
		Preload("Ideas.Books.Authors").
		Preload("Ideas.Movies").
		Preload("Ideas.Movies.Genres").
		Preload("Ideas.Songs").
		Preload("Ideas.Songs.MusicalAlbum").
		Preload("Ideas.Songs.Artists").
		Preload("Ideas.Wikis").
		Preload("Ideas.Generics").
		Preload("Ideas.Articles").
		First(&tag, tagReq.ID)

	return tag, nil
}

func (ir *IdeaRepository) FetchGraph(c context.Context, ID string) (domain.Idea, error) {
	var idea domain.Idea

	res := ir.database.
		Preload("Blogs.Ideas").
		Preload("Videos.Ideas").
		Preload("Anchors.Anchors").
		Preload("Books.Authors").
		Preload("Movies.Ideas").
		Preload("Songs.Ideas").
		Preload("Wikis.Ideas").
		Preload("Generics.Ideas").
		Preload("Articles.Ideas").
		First(&idea, ID)

	return idea, res.Error
}

func (ir *IdeaRepository) FetchByResourceID(c context.Context, resType string, resID string) []domain.Idea {
	var ideas []domain.Idea

	switch resType {
	case "blogs":
		blog := domain.Blog{}
		ir.database.Preload("Ideas").First(&blog, resID)
		ideas = blog.Ideas
	case "videos":
		video := domain.Video{}
		ir.database.Preload("Ideas").First(&video, resID)
		ideas = video.Ideas
	case "books":
		book := domain.Book{}
		ir.database.Preload("Ideas").First(&book, resID)
		ideas = book.Ideas
	case "movies":
		movie := domain.Movie{}
		ir.database.Preload("Ideas").First(&movie, resID)
		ideas = movie.Ideas
	case "songs":
		song := domain.Song{}
		ir.database.Preload("Ideas").First(&song, resID)
		ideas = song.Ideas
	case "wikis":
		wiki := domain.Wiki{}
		ir.database.Preload("Ideas").First(&wiki, resID)
		ideas = wiki.Ideas
	case "generics":
		generic := domain.Generic{}
		ir.database.Preload("Ideas").First(&generic, resID)
		ideas = generic.Ideas
	case "articles":
		article := domain.Article{}
		ir.database.Preload("Ideas").First(&article, resID)
		ideas = article.Ideas
	case "anchors":
		anchor := domain.Idea{}
		ir.database.Preload("Anchors").First(&anchor, resID)
		ideas = anchor.Anchors
	}

	return ideas
}

func (ir *IdeaRepository) DeleteByID(c context.Context, id string) error {
	var idea domain.Idea
	tx := ir.database.Delete(&idea, id)
	return tx.Error
}

// beforeCreate hook?
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

	for iw, wiki := range idea.Wikis {
		if wiki.ID == 0 {
			w := domain.Wiki{}
			ir.database.Where(&domain.Wiki{Url: wiki.Url}).First(&w)

			if w.ID != 0 {
				wikiPtr := &idea.Wikis[iw]
				wikiPtr.ID = w.ID
			}
		}
	}

	for ig, generic := range idea.Generics {
		if generic.ID == 0 {
			g := domain.Generic{}
			ir.database.Where(&domain.Generic{Url: generic.Url}).First(&g)

			if g.ID != 0 {
				genericPtr := &idea.Generics[ig]
				genericPtr.ID = g.ID
			}
		}
	}

	for ai, article := range idea.Articles {
		if article.ID == 0 {
			a := domain.Article{}
			ir.database.Where(&domain.Article{Url: article.Url}).First(&a)

			if a.ID != 0 {
				articlePtr := &idea.Articles[ai]
				articlePtr.ID = a.ID
			}
		}
	}
}

// beforeCreate hook?
func (ir *IdeaRepository) assignResourceFields(idea *domain.Idea) {
	// assign youtube channel as well coming from the frontend request fields
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
