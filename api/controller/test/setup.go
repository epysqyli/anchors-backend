package controller

import (
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

func cleanupUser(db *gorm.DB, userName string) {
	var user domain.User
	db.Model(&domain.User{}).Where("name = ?", userName).First(&user)
	db.Unscoped().Delete(&user, "name = ?", userName)
}
