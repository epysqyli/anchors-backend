package domain

import (
	"context"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Task     []Task
}

type UserRepository interface {
	GetByID(c context.Context, id string) (User, error)
	GetByEmail(c context.Context, email string) (User, error)
	Create(c context.Context, user *User) error
	Fetch(c context.Context) ([]User, error)
}
