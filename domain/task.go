package domain

import (
	"context"
)

type Task struct {
	UserId uint   `json:"user_id" form:"user_id"`
	Title  string `json:"title" form:"title" binding:"required"`
}

type TaskRepository interface {
	Create(c context.Context, task *Task) error
	FetchByUserID(c context.Context, userID string) ([]Task, error)
}

type TaskUsecase interface {
	Create(c context.Context, task *Task) error
	FetchByUserID(c context.Context, userID string) ([]Task, error)
}
