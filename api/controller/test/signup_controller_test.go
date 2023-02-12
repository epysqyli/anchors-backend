package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	routev1 "github.com/epysqyli/anchors-backend/api/route/v1"
	"github.com/epysqyli/anchors-backend/bootstrap"
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// does it make sense to have controller_test as a package?
// can it be done in TestMain for the controller package?
func setup() (*gin.Engine, *gorm.DB) {
	app := bootstrap.App("../../../.env")
	psqlDB := app.Postgres

	gin.SetMode(gin.TestMode)
	gin := gin.New()
	routerV1 := gin.Group("v1")
	routev1.Setup(app.Env, psqlDB, routerV1)

	return gin, psqlDB
}

func cleanup(db *gorm.DB, userName string) {
	var user domain.User
	db.Model(&domain.User{}).Where("name = ?", userName).First(&user)
	db.Unscoped().Delete(&user, "name = ?", userName)
}

func TestSignup(t *testing.T) {
	gin, db := setup()

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

	cleanup(db, "anchors")
}
