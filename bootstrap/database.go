package bootstrap

import (
	"fmt"

	"github.com/epysqyli/anchors-backend/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDatabase(env *Env) *gorm.DB {
	connString := fmt.Sprintf("host=postgres user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/Rome", env.PostgresUser, env.PostgresPassword, env.PostgresDB)

	config := postgres.Config{
		DSN:                  connString,
		PreferSimpleProtocol: true,
	}

	db, err := gorm.Open(postgres.New(config), &gorm.Config{})

	if err != nil {
		panic("database connection failed")
	}

	db.AutoMigrate(&domain.User{}, &domain.Task{})

	return db
}
