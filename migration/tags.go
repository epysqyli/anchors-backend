package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Tags() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "19",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.Tag{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("tags")
		},
	}
}
