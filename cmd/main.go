package main

import (
	"time"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/api/middleware"
	routeV1 "github.com/amitshekhariitbhu/go-backend-clean-architecture/api/route/v1"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/bootstrap"
	"github.com/gin-gonic/gin"
)

func main() {
	app := bootstrap.App()
	env := app.Env
	psqlDB := app.Postgres
	timeout := time.Duration(env.ContextTimeout) * time.Second
	gin := gin.Default()
	gin.Use(middleware.Cors())
	routerV1 := gin.Group("v1")
	routeV1.Setup(env, timeout, psqlDB, routerV1)

	gin.Run(env.ServerAddress)
}
