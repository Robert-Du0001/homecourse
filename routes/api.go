package routes

import (
	"github.com/goravel/framework/contracts/route"

	"homecourse/app/facades"
	"homecourse/app/http/controllers"
	"homecourse/app/http/middleware"
)

func Api() {
	userController := controllers.NewUserController()
	categoryController := controllers.NewCategoryController()
	courseController := controllers.NewCourseController()
	episodeController := controllers.NewEpisodeController()

	facades.Route().Prefix("api").Group(func(router route.Router) {
		router.Post("user", userController.Store)
		router.Post("user/token", userController.Login)

		router.Middleware(middleware.Auth()).Group(func(router route.Router) {
			router.Get("user", userController.Show)

			// 课程分类相关路由
			router.Get("categories", categoryController.Index)

			// 课程相关路由
			router.Get("courses", courseController.Index)
			router.Get("courses/{id}", courseController.Show)

			// 剧集相关路由
			router.Get("episodes", episodeController.Index)
			router.Get("episodes/{id}", episodeController.Show)

			// 管理员
			router.Middleware(middleware.Admin()).Prefix("admin").Group(func(router route.Router) {
				router.Get("users", userController.Index)
				router.Delete("users/{id}", userController.Destroy)

				// 课程分类相关路由
				router.Post("categories", categoryController.Store)
				router.Put("categories/{id}", categoryController.Update)
				router.Put("categories/{id}/default", categoryController.UpdateDefault)
				router.Put("categories/sort", categoryController.UpdateSort)
				router.Delete("categories/{id}", categoryController.Destroy)

				// 课程相关路由
				router.Get("courses", courseController.AdminIndex)
				router.Post("courses", courseController.Store)
				router.Put("courses/{id}", courseController.Update)
				router.Put("courses/sort", courseController.UpdateSort)
				router.Delete("courses/{id}", courseController.Destroy)

				// 剧集相关路由
				router.Put("episodes/scan", episodeController.Scan)
			})
		})
	})
}
