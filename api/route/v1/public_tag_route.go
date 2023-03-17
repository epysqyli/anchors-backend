package route

import (
	"github.com/epysqyli/anchors-backend/api/controller"
	"github.com/epysqyli/anchors-backend/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewPublicTagRouter(db *gorm.DB, group *gin.RouterGroup) {
	tc := &controller.TagsController{
		TagRepository: repository.NewTagRepository(db),
	}

	group.GET("/tags", tc.FetchAllTags)
	group.GET("/tags/:id", tc.FetchByID)
	group.GET("/tags/name/:name", tc.FetchByName)
}
