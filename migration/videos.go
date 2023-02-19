package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Videos() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "3",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.Video{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("videos")
		},
	}
}
