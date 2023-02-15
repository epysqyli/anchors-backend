package controller

import (
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/gin-gonic/gin"
)

type IdeaController struct {
	IdeaRepository domain.IdeaRepository
}

func (ic *IdeaController) FetchIdeaByID(ctx *gin.Context) {}

func (ic *IdeaController) FetchIdeasByUserID(ctx *gin.Context) {}

func (ic *IdeaController) FetchAllIdeas(ctx *gin.Context) {}

func (ic *IdeaController) CreateIdea(ctx *gin.Context) {}

func (ic *IdeaController) DeleteIdeaByID(ctx *gin.Context) {}
