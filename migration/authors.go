package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Authors() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "8",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.Author{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("authors")
		},
	}
}
