package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	routev1 "github.com/epysqyli/anchors-backend/api/route/v1"
	"github.com/epysqyli/anchors-backend/bootstrap"
	"github.com/gin-gonic/gin"
)

func TestSignup(t *testing.T) {
	app := bootstrap.App("../../../.env")
	psqlDB := app.Postgres

	gin.SetMode(gin.TestMode)
	gin := gin.Default()
	routerV1 := gin.Group("v1")
	routev1.Setup(app.Env, psqlDB, routerV1)

	t.Run("success", func(t *testing.T) {
		body := []byte(`{
			"name": "anchors",
			"email": "anchors@gmail.com",
			"password": "anchors"
		}`)

		bodyReader := bytes.NewReader(body)

		req, err := http.NewRequest(http.MethodPost, "/v1/signup", bodyReader)

		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		req.Header.Add("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		gin.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("Response returned with an unexpected status: %v\n", rec.Code)
		}
	})
}
