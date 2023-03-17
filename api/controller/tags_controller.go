package controller

import (
	"fmt"
	"net/http"

	"github.com/epysqyli/anchors-backend/domain"
	"github.com/gin-gonic/gin"
)

type TagsController struct {
	TagRepository domain.TagRepository
}

func (tc *TagsController) FetchAllTags(ctx *gin.Context) {
	tags := tc.TagRepository.FetchAll()

	ctx.JSON(http.StatusOK, tags)
}

func (tc *TagsController) CreateTag(ctx *gin.Context) {
	tag := domain.Tag{}
	ctx.ShouldBind(&tag)

	err := tc.TagRepository.Create(&tag)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Could not create tag"})
	}

	ctx.JSON(http.StatusCreated, tag)
}

func (tc *TagsController) FetchByID(ctx *gin.Context) {
	ID := ctx.Param("id")
	tag := tc.TagRepository.FetchById(ID)

	ctx.JSON(http.StatusOK, tag)
}

func (tc *TagsController) FetchByName(ctx *gin.Context) {
	name := ctx.Param("name")
	tag := tc.TagRepository.FetchByName(name)

	ctx.JSON(http.StatusOK, tag)
}

func (tc *TagsController) DeleteByID(ctx *gin.Context) {
	ID := ctx.Param("id")
	err := tc.TagRepository.DeleteByID(ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Could not delete tag"})
	}

	ctx.JSON(http.StatusAccepted, domain.SuccessResponse{Message: fmt.Sprintf("Tag with id %s deleted", ID)})
}
