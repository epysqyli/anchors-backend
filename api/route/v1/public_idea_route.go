package route

import (
	"github.com/epysqyli/anchors-backend/api/controller"
	"github.com/epysqyli/anchors-backend/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewPublicIdeaRouter(psqlDB *gorm.DB, group *gin.RouterGroup) {
	ic := &controller.IdeaController{
		IdeaRepository: repository.NewIdeaRepository(psqlDB),
	}

	group.GET("/ideas", ic.FetchAllIdeas)
	group.GET("/ideas/:id", ic.FetchIdeaByID)
	group.GET("/users/:user_id/ideas", ic.FetchIdeasByUserID)
	group.GET("/ideas/graph/:id", ic.FetchGraphByIdeaID)
	group.GET("/ideas/by_anchor/:resource_type/:resource_id", ic.FetchIdeaByResourceID)
	group.GET("/ideas/tags", ic.FetchByTags)
}
