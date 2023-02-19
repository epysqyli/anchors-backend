package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Blogs() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "5",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.Blog{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("blogs")
		},
	}
}
