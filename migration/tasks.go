package migration

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Tasks() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "2",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.Task{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("tasks")
		},
	}
}
