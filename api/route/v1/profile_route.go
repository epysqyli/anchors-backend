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

func NewProfileRouter(env *bootstrap.Env, timeout time.Duration, psqlDB *gorm.DB, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(psqlDB)
	pc := &controller.ProfileController{
		ProfileUsecase: usecase.NewProfileUsecase(ur, timeout),
	}
	group.GET("/profile", pc.Fetch)
}
