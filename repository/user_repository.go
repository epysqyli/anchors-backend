package repository

import (
	"context"
	"strconv"
	"strings"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	database *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{database: db}
}

func (ur *userRepository) GetByID(c context.Context, id string) (domain.User, error) {
	var user domain.User
	userId, _ := strconv.ParseInt(strings.Split(id, ".")[0], 10, 0)
	res := ur.database.First(&user, userId)
	return user, res.Error
}

func (ur *userRepository) GetByEmail(c context.Context, email string) (domain.User, error) {
	var user domain.User
	res := ur.database.Model(&domain.User{}).Where("email = ?", email).First(&user)

	return user, res.Error
}

func (ur *userRepository) Create(c context.Context, user *domain.User) error {
	tx := ur.database.Create(user)
	return tx.Error
}

func (ur *userRepository) Fetch(c context.Context) ([]domain.User, error) {
	var users []domain.User
	return users, nil
}
