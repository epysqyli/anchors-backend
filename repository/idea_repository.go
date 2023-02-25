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
	res := ir.database.Create(idea)
	if res.Error != nil {
		return res.Error
	}

	ir.assignRelationFields(idea)
	return nil
}

func (ir *IdeaRepository) FetchAll(c context.Context) ([]domain.Idea, error) {
	var ideas []domain.Idea
	res := ir.database.Find(&ideas)
	return ideas, res.Error
}

func (ir *IdeaRepository) FetchByUserID(c context.Context, userId string) ([]domain.Idea, error) {
	var ideas []domain.Idea
	res := ir.database.Model(&domain.Idea{}).
		Preload("Blogs").
		Preload("Videos").
		Preload("Anchors").
		Find(&ideas, "user_id = ?", userId)

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
		First(&idea, id)

	return idea, res.Error
}

func (ir *IdeaRepository) DeleteByID(c context.Context, id string) error {
	var idea domain.Idea
	tx := ir.database.Delete(&idea, id)
	return tx.Error
}

/**
 *	can this be done with an afterCreate DB hook?
 *  can db calls be limited to a single call?
 */
func (ir *IdeaRepository) assignRelationFields(idea *domain.Idea) {
	for _, video := range idea.Videos {
		ir.database.
			Model(domain.IdeasVideos{IdeaID: idea.ID, VideoID: video.ID}).
			Update("timestamp", video.Timestamp)
	}
}
