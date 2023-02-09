package controller

import "github.com/epysqyli/anchors-backend/domain"

type IdeaController struct {
	IdeaRepository domain.IdeaRepository
}

// setup test for this?
// 1 - call to endpoint
// 2 - controller action
// 3 - database actions via repository
// 4 - json response
