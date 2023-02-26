package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Books() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "7",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.Book{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("books")
		},
	}
}
