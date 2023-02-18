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
	return res.Error
}

func (ir *IdeaRepository) FetchAll(c context.Context) ([]domain.Idea, error) {
	var ideas []domain.Idea
	res := ir.database.Find(&domain.Idea{})
	return ideas, res.Error
}

func (ir *IdeaRepository) FetchByUserID(c context.Context, userId string) ([]domain.Idea, error) {
	var ideas []domain.Idea
	res := ir.database.Model(&domain.Idea{}).Find(&ideas, "user_id = ?", userId)
	return ideas, res.Error
}
