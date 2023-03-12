package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Generics() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "17",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.Generic{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("generics")
		},
	}
}
