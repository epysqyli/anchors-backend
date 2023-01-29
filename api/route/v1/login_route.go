package route

import (
	"time"

	"github.com/epysqyli/anchors-backend/api/controller"
	"github.com/epysqyli/anchors-backend/bootstrap"
	"github.com/epysqyli/anchors-backend/repository"
	"github.com/epysqyli/anchors-backend/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewLoginRouter(env *bootstrap.Env, timeout time.Duration, psqlDB *gorm.DB, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(psqlDB)
	lc := &controller.LoginController{
		LoginUsecase: usecase.NewLoginUsecase(ur, timeout),
		Env:          env,
	}
	group.POST("/login", lc.Login)
}
