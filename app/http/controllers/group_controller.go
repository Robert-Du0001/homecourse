package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type GroupController struct {
	// Dependent services
}

func NewGroupController() *GroupController {
	return &GroupController{
		// Inject services
	}
}

func (r *GroupController) Index(ctx http.Context) http.Response {
	return nil
}
