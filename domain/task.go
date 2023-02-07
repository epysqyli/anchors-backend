package domain

type Task struct {
	UserId uint   `json:"user_id" form:"user_id"`
	Title  string `json:"title" form:"title" binding:"required"`
}

type TaskRepository interface {
	Create(task *Task) error
	FetchByUserID(userID string) ([]Task, error)
}
