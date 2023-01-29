package repository

import (
	"context"

	"github.com/epysqyli/anchors-backend/domain"
	"gorm.io/gorm"
)

type taskRepository struct {
	database *gorm.DB
}

func NewTaskRepository(db *gorm.DB) domain.TaskRepository {
	return &taskRepository{
		database: db,
	}
}

func (tr *taskRepository) Create(c context.Context, task *domain.Task) error {
	// insert task into DB
	// -> deal with context as well (timeout?)

	res := tr.database.Create(task)
	return res.Error
}

func (tr *taskRepository) FetchByUserID(c context.Context, userID string) ([]domain.Task, error) {
	// fetch all user's tasks
	// -> deal with context as well (timeout?)

	var tasks []domain.Task
	res := tr.database.Model(&domain.Task{}).Find(&tasks, "user_id = ?", userID)
	return tasks, res.Error
}
