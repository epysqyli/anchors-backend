package controller

import (
	"fmt"
	"log"
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

func (ic *IdeaController) FetchGraphByIdeaID(ctx *gin.Context) {
	id := ctx.Param("id")

	idea, err := ic.IdeaRepository.FetchGraph(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, idea)
}

func (ic *IdeaController) FetchIdeaByResourceID(ctx *gin.Context) {
	resType := ctx.Param("resource_type")
	resID := ctx.Param("resource_id")
	ideas := ic.IdeaRepository.FetchByResourceID(ctx, resType, resID)

	ctx.JSON(http.StatusOK, ideas)
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

/**
 * tag ids are passed as hyphen delimited strings
 * /v1/ideas/tags?and=1-2-3-4-5
 * these are captured and processed accordingly in order
 * to build the TagIdeasRequest struct used in the repo
 */
func (ic *IdeaController) FetchByTags(ctx *gin.Context) {
	tagQuery := domain.NewTagQuery(ctx)
	log.Default().Printf("%+v", tagQuery)

	tag, err := ic.IdeaRepository.FetchByTags(tagQuery)

	if err != nil {
		ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tag)
}

func (ic *IdeaController) FetchAllIdeas(ctx *gin.Context) {
	ideas, err := ic.IdeaRepository.FetchAll(ctx)

	if err != nil {
		ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ideas)
}

func (ic *IdeaController) CreateIdea(ctx *gin.Context) {
	var idea domain.Idea

	err := ctx.ShouldBind(&idea)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	if idea.HasNoResources() {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "An idea needs at least one type of resource associated to it",
		})

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
