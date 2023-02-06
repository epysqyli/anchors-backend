package usecase

import (
	"context"
	"time"

	"github.com/epysqyli/anchors-backend/domain"
)

type ideaUsecase struct {
	ideaRepository domain.IdeaRepository
	contextTimeout time.Duration
}

func NewIdeaUsecase(ideaRepository domain.IdeaRepository, timeout time.Duration) domain.IdeaUsecase {
	return &ideaUsecase{
		ideaRepository: ideaRepository,
		contextTimeout: timeout,
	}
}

func (iu *ideaUsecase) Create(c context.Context, idea *domain.Idea) error {
	ctx, cancel := context.WithTimeout(c, iu.contextTimeout)
	defer cancel()
	return iu.ideaRepository.Create(ctx, idea)
}

func (iu *ideaUsecase) FetchByUserID(c context.Context, userID string) ([]domain.Idea, error) {
	ctx, cancel := context.WithTimeout(c, iu.contextTimeout)
	defer cancel()
	return iu.ideaRepository.FetchByUserID(ctx, userID)
}
