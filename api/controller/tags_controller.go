package controller

import (
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
