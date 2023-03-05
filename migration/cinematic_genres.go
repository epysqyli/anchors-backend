package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CinematicGenres() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "10",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.CinematicGenre{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("cinematic_genres")
		},
	}
}
