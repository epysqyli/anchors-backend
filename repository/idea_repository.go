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

// the main general feed: pagination needs to be set up
func (ir *IdeaRepository) FetchAll(c context.Context) ([]domain.Idea, error) {
	var ideas []domain.Idea
	res := ir.database.
		Preload("Blogs").
		Preload("Videos").
		Preload("Anchors").
		Preload("Books.Authors").
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
		Find(&ideas, "user_id = ?", userID)

	return ideas, res.Error
}

/**
 * each 'preload' executes a query
 * can gorm 'joins' be used to execute a single query on many to many tables?
 * optimize when and if necessary
 */
func (ir *IdeaRepository) FetchByID(c context.Context, id string) (domain.Idea, error) {
	var idea domain.Idea
	res := ir.database.
		Preload("Blogs").
		Preload("Videos").
		Preload("Anchors").
		Preload("Books.Authors").
		First(&idea, id)

	return idea, res.Error
}

func (ir *IdeaRepository) DeleteByID(c context.Context, id string) error {
	var idea domain.Idea
	tx := ir.database.Delete(&idea, id)
	return tx.Error
}

/**
 * assign IDs to existing resources
 * needed for all those resource provided by external APIs:
 * books + authors, movies, songs
 * these resources are not unique based on their URL
 */
func (ir *IdeaRepository) assignExistingIDs(idea *domain.Idea) {
	if idea.Books == nil {
		return
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
}

/**
 * implement logic to assign unique identifier for each resource type
 * example: youtube video has unique ID from youtube, other videos should use URLs
 * this should be extended to other resource types
 * it maybe can be refactored into model methods if access to DB is postponed
 */
func (ir *IdeaRepository) assignResourceFields(idea *domain.Idea) {
	for i := range idea.Videos {
		videaPtr := &idea.Videos[i]
		videaPtr.AssignIdentifier()
	}
}

/**
 *	can this be done with an afterCreate DB hook?
 *  can db calls be limited to a single call?
 */
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
}
