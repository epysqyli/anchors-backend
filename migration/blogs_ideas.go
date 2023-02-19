package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func BlogsIdeas() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "6",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.IdeasVideos{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("blogs_ideas")
		},
	}
}
