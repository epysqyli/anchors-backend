package controller

import (
	"testing"

	routeV1 "github.com/epysqyli/anchors-backend/api/route/v1"
	"github.com/epysqyli/anchors-backend/bootstrap"
	"github.com/gin-gonic/gin"
)

func TestCreateTask(t *testing.T) {
	app := bootstrap.App("../../../.env")
	psqlDB := app.Postgres

	gin.SetMode(gin.TestMode)
	gin := gin.Default()
	routerV1 := gin.Group("v1")
	routeV1.Setup(app.Env, psqlDB, routerV1)

	// run auth method -> will have to be handled differently for google auth (maybe both auth modes can be kept?)
	// assign bearer to req header
	// perform request to /v1/tasks
	// inspect response

	t.Run("success", func(t *testing.T) {
		// req, err := http.NewRequest(http.MethodPost, "/v1/tasks", nil)

	})
}
