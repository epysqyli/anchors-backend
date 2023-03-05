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
	port := ""

	switch env.AppEnv {
	case "development":
		dbName = env.PostgresDevelopmentDB
		host = "postgres"
		user = env.PostgresDevelopmentUser
		password = env.PostgresDevelopmentPassword
		port = env.DBPortHostDevelopment
	case "test":
		dbName = env.PostgresTestDB
		host = "localhost"
		user = env.PostgresTestUser
		password = env.PostgresTestPassword
		port = env.DBPortHostTest
	case "production":
		dbName = env.PostgresProductionDB
		host = ""
		user = env.PostgresProductionUser
		password = env.PostgresProductionPassword
		port = ""
	}

	connString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Rome",
		host, user, password, dbName, port,
	)

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
		migration.Ideas(),
		migration.Videos(),
		migration.IdeasVideos(),
		migration.Blogs(),
		migration.BlogsIdeas(),
		migration.Books(),
		migration.Authors(),
		migration.BooksIdeas(),
		migration.CinematicGenres(),
		migration.Movies(),
	})

	if err = migrations.Migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	return db
}
