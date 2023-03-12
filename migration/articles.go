package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Articles() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "18",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.Article{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("articles")
		},
	}
}
