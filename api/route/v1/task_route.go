package route

import (
	"github.com/epysqyli/anchors-backend/api/controller"
	"github.com/epysqyli/anchors-backend/bootstrap"
	"github.com/epysqyli/anchors-backend/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewTaskRouter(env *bootstrap.Env, psqlDB *gorm.DB, group *gin.RouterGroup) {
	tc := &controller.TaskController{
		TaskRepository: repository.NewTaskRepository(psqlDB),
	}

	group.GET("/tasks", tc.FetchTask)
	group.POST("/tasks", tc.CreateTask)
}
