package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func BooksIdeas() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "9",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.BooksIdeas{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("books_ideas")
		},
	}
}
