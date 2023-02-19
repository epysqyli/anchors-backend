package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Resources() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "3",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.Resource{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("resources")
		},
	}
}
