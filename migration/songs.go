package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Songs() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "14",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.Song{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("songs")
		},
	}
}
