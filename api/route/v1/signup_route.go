package route

import (
	"github.com/epysqyli/anchors-backend/api/controller"
	"github.com/epysqyli/anchors-backend/bootstrap"
	"github.com/epysqyli/anchors-backend/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewSignupRouter(env *bootstrap.Env, psqlDB *gorm.DB, group *gin.RouterGroup) {
	sc := controller.SignupController{
		UserRepository: repository.NewUserRepository(psqlDB),
		Env:            env,
	}

	group.POST("/signup", sc.Signup)
}
