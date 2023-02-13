package bootstrap

import (
	"gorm.io/gorm"
)

type Application struct {
	Env      *Env
	Postgres *gorm.DB
}

func App(envPath string, envMode string) Application {
	app := &Application{}
	app.Env = NewEnv(envPath, envMode)
	app.Postgres = NewPostgresDatabase(app.Env)
	return *app
}
