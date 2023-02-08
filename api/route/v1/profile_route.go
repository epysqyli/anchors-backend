package route

import (
	"github.com/epysqyli/anchors-backend/api/controller"
	"github.com/epysqyli/anchors-backend/bootstrap"
	"github.com/epysqyli/anchors-backend/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewProfileRouter(env *bootstrap.Env, psqlDB *gorm.DB, group *gin.RouterGroup) {
	pc := &controller.ProfileController{
		UserRepository: repository.NewUserRepository(psqlDB),
	}

	group.GET("/profile", pc.Fetch)
}
