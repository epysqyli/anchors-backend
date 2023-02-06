package bootstrap

import (
	"fmt"
	"log"

	"github.com/epysqyli/anchors-backend/migration"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDatabase(env *Env) *gorm.DB {
	connString := fmt.Sprintf("host=postgres"+
		" user=%s"+
		" password=%s"+
		" dbname=%s"+
		" port=5432"+
		" sslmode=disable"+
		" TimeZone=Europe/Rome",
		env.PostgresUser,
		env.PostgresPassword,
		env.PostgresDB)

	config := postgres.Config{
		DSN:                  connString,
		PreferSimpleProtocol: true,
	}

	db, err := gorm.Open(postgres.New(config), &gorm.Config{})

	if err != nil {
		panic("database connection failed")
	}

	migrations := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migration.Users(),
		migration.Tasks(),
		migration.Ideas(),
		migration.Resources(),
		migration.IdeasResources(),
	})

	if err = migrations.Migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	return db
}
