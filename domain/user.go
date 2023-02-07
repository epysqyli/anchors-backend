package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Task     []Task
}

type UserRepository interface {
	GetByID(id string) (User, error)
	GetByEmail(email string) (User, error)
	Create(user *User) error
	Fetch() ([]User, error)
}
