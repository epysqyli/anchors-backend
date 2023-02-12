package bootstrap

import (
	"gorm.io/gorm"
)

type Application struct {
	Env      *Env
	Postgres *gorm.DB
}

func App(envPath string) Application {
	app := &Application{}
	app.Env = NewEnv(envPath)
	app.Postgres = NewPostgresDatabase(app.Env)
	return *app
}
