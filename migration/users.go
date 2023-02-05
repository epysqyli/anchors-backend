package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Users() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.User{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("users")
		},
	}
}
