package route

import (
	"github.com/epysqyli/anchors-backend/api/controller"
	"github.com/epysqyli/anchors-backend/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewProtectedIdeaRouter(psqlDB *gorm.DB, group *gin.RouterGroup) {
	ic := &controller.IdeaController{
		IdeaRepository: repository.NewIdeaRepository(psqlDB),
	}

	group.POST("/ideas", ic.CreateIdea)
	group.DELETE("/ideas/:id", ic.DeleteIdeaByID)
}
