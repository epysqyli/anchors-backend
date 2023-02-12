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
	dbName := ""
	host := ""
	user := ""
	password := ""

	switch env.AppEnv {
	case "development":
		dbName = env.PostgresDevelopmentDB
		host = "postgres"
		user = env.PostgresUser
		password = ""
	case "test":
		dbName = env.PostgresTestDB
		host = "localhost"
		user = env.PostgresTestUser
		password = env.PostgresTestPassword
	case "production":
		dbName = ""
		host = ""
		user = ""
		password = env.PostgresPassword
	}

	connString := fmt.Sprintf("host=%s"+
		" user=%s"+
		" password=%s"+
		" dbname=%s"+
		" port=5432"+
		" sslmode=disable"+
		" TimeZone=Europe/Rome",
		host, user, password, dbName)

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
