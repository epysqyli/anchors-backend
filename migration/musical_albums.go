package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func MusicalAlbums() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "13",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.MusicalAlbum{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("musical_albums")
		},
	}
}
