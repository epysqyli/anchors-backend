package route

import (
	"time"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/api/middleware"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/bootstrap"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(env *bootstrap.Env, timeout time.Duration, psqlDB *gorm.DB, routerV1 *gin.RouterGroup) {
	publicRouterV1 := routerV1.Group("")
	// All Public APIs
	NewSignupRouter(env, timeout, psqlDB, publicRouterV1)
	NewLoginRouter(env, timeout, psqlDB, publicRouterV1)
	NewRefreshTokenRouter(env, timeout, psqlDB, publicRouterV1)

	protectedRouterV1 := routerV1.Group("")
	// Middleware to verify AccessToken
	protectedRouterV1.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	// All Private APIs
	NewProfileRouter(env, timeout, psqlDB, protectedRouterV1)
	NewTaskRouter(env, timeout, psqlDB, protectedRouterV1)
}
