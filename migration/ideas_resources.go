package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func IdeasResources() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "5",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.IdeaResource{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("ideas_resources")
		},
	}
}
