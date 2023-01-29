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

func NewTaskRouter(env *bootstrap.Env, timeout time.Duration, psqlDB *gorm.DB, group *gin.RouterGroup) {
	tr := repository.NewTaskRepository(psqlDB)
	tc := &controller.TaskController{
		TaskUsecase: usecase.NewTaskUsecase(tr, timeout),
	}

	group.GET("/tasks", tc.Fetch)
	group.POST("/tasks", tc.Create)
}
