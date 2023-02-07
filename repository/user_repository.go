package repository

import (
	"strconv"
	"strings"

	"github.com/epysqyli/anchors-backend/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	database *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{database: db}
}

func (ur *userRepository) GetByID(id string) (domain.User, error) {
	var user domain.User
	userId, _ := strconv.ParseInt(strings.Split(id, ".")[0], 10, 0)
	res := ur.database.First(&user, userId)
	return user, res.Error
}

func (ur *userRepository) GetByEmail(email string) (domain.User, error) {
	var user domain.User
	res := ur.database.Model(&domain.User{}).Where("email = ?", email).First(&user)

	return user, res.Error
}

func (ur *userRepository) Create(user *domain.User) error {
	tx := ur.database.Create(user)
	return tx.Error
}

func (ur *userRepository) Fetch() ([]domain.User, error) {
	var users []domain.User
	return users, nil
}
