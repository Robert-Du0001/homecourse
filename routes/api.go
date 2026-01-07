package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"homecourse/app/http/controllers"
)

func Api() {
	userController := controllers.NewUserController()

	facades.Route().Prefix("api").Group(func(router route.Router) {
		router.Get("/users/{id}", userController.Show)
	})
}
