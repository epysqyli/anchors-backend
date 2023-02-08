package route

import (
	"github.com/epysqyli/anchors-backend/api/controller"
	"github.com/epysqyli/anchors-backend/bootstrap"
	"github.com/epysqyli/anchors-backend/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRefreshTokenRouter(env *bootstrap.Env, psqlDB *gorm.DB, group *gin.RouterGroup) {
	rtc := &controller.RefreshTokenController{
		UserRepository: repository.NewUserRepository(psqlDB),
		Env:            env,
	}

	group.POST("/refresh", rtc.RefreshToken)
}
