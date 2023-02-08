package route

import (
	"github.com/epysqyli/anchors-backend/api/controller"
	"github.com/epysqyli/anchors-backend/bootstrap"
	"github.com/epysqyli/anchors-backend/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewLoginRouter(env *bootstrap.Env, psqlDB *gorm.DB, group *gin.RouterGroup) {
	lc := &controller.LoginController{
		UserRepository: repository.NewUserRepository(psqlDB),
		Env:            env,
	}

	group.POST("/login", lc.Login)
}
