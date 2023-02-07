package controller

import (
	"net/http"

	"github.com/epysqyli/anchors-backend/domain"
	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	UserRepository domain.UserRepository
}

func (pc *ProfileController) Fetch(c *gin.Context) {
	userID := c.GetString("x-user-id")

	user, err := pc.UserRepository.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	profile := domain.Profile{
		Name:  user.Name,
		Email: user.Email,
	}

	c.JSON(http.StatusOK, profile)
}
