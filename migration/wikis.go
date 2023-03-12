package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Wikis() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "16",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.Wiki{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("wikis")
		},
	}
}
