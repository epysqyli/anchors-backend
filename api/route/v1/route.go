package route

import (
	"github.com/epysqyli/anchors-backend/api/middleware"
	"github.com/epysqyli/anchors-backend/bootstrap"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(env *bootstrap.Env, psqlDB *gorm.DB, routerV1 *gin.RouterGroup) {
	publicRouterV1 := routerV1.Group("")
	NewSignupRouter(env, psqlDB, publicRouterV1)
	NewLoginRouter(env, psqlDB, publicRouterV1)
	NewRefreshTokenRouter(env, psqlDB, publicRouterV1)
	NewPublicIdeaRouter(psqlDB, publicRouterV1)
	NewPublicTagRouter(psqlDB, publicRouterV1)

	protectedRouterV1 := routerV1.Group("")
	protectedRouterV1.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	NewProfileRouter(env, psqlDB, protectedRouterV1)
	NewProtectedIdeaRouter(psqlDB, protectedRouterV1)
}
