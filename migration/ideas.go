package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Ideas() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "2",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.Idea{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("ideas")
		},
	}
}
