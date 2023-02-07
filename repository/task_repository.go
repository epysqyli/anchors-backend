package repository

import (
	"github.com/epysqyli/anchors-backend/domain"
	"gorm.io/gorm"
)

type TaskRepository struct {
	database *gorm.DB
}

func NewTaskRepository(db *gorm.DB) domain.TaskRepository {
	return &TaskRepository{
		database: db,
	}
}

func (tr *TaskRepository) Create(task *domain.Task) error {
	res := tr.database.Create(task)
	return res.Error
}

func (tr *TaskRepository) FetchByUserID(userID string) ([]domain.Task, error) {
	var tasks []domain.Task
	res := tr.database.Model(&domain.Task{}).Find(&tasks, "user_id = ?", userID)
	return tasks, res.Error
}
