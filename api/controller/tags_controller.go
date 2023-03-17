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
