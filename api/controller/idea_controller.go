package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/epysqyli/anchors-backend/domain"
	"github.com/gin-gonic/gin"
)

type IdeaController struct {
	IdeaRepository domain.IdeaRepository
}

func (ic *IdeaController) FetchIdeaByID(ctx *gin.Context) {
	id := ctx.Param("id")
	idea, err := ic.IdeaRepository.FetchByID(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, idea)
}

func (ic *IdeaController) FetchIdeasByUserID(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	ideas, err := ic.IdeaRepository.FetchByUserID(ctx, userID)

	if err != nil {
		ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ideas)
}

func (ic *IdeaController) FetchAllIdeas(ctx *gin.Context) {
	ideas, err := ic.IdeaRepository.FetchAll(ctx)

	if err != nil {
		ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ideas)
}

/*
* request comes in
* associated resources are created based on the resource type
* JSON resp is returned with all the necessary associations
 */
func (ic *IdeaController) CreateIdea(ctx *gin.Context) {
	var idea domain.Idea

	err := ctx.ShouldBind(&idea)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	userID, _ := strconv.ParseInt(ctx.GetString("x-user-id"), 0, 32)
	idea.UserID = uint(userID)

	err = ic.IdeaRepository.Create(ctx, &idea)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, idea)
}

func (ic *IdeaController) DeleteIdeaByID(ctx *gin.Context) {
	id := ctx.Param("id")
	err := ic.IdeaRepository.DeleteByID(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, domain.SuccessResponse{Message: fmt.Sprintf("Idea with id %s deleted", id)})
}
