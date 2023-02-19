package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func IdeasVideos() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "4",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.IdeasVideos{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("ideas_videos")
		},
	}
}
