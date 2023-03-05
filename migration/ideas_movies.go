package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func IdeasMovies() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "12",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.IdeasMovies{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("ideas_movies")
		},
	}
}
