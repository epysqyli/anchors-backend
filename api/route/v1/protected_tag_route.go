package route

import (
	"github.com/epysqyli/anchors-backend/api/controller"
	"github.com/epysqyli/anchors-backend/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewProtectedTagRouter(db *gorm.DB, group *gin.RouterGroup) {
	tc := &controller.TagsController{
		TagRepository: repository.NewTagRepository(db),
	}

	group.POST("/tags", tc.CreateTag)
}
