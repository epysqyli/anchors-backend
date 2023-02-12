package controller

import (
	"net/http"
	"strconv"

	"github.com/epysqyli/anchors-backend/domain"
	"github.com/gin-gonic/gin"
)

type TaskController struct {
	TaskRepository domain.TaskRepository
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	var task domain.Task

	err := c.ShouldBind(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	userId, _ := strconv.ParseInt(c.GetString("x-user-id"), 0, 32)
	task.UserId = uint(userId)

	err = tc.TaskRepository.Create(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Task created successfully",
	})
}

func (u *TaskController) FetchTask(c *gin.Context) {
	userID := c.GetString("x-user-id")

	tasks, err := u.TaskRepository.FetchByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
