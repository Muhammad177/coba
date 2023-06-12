package routes

import (
	"Capstone/constant"
	"Capstone/controller"
	"Capstone/midleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRoutes(e *echo.Echo) {

	midleware.LogMiddleware(e)
	// routing with query parameter
	e.POST("/login", controller.LoginController)
	e.POST("/login/admin", controller.LoginAdminController)
	e.POST("/user", controller.CreateUserController)

	eJwt := e.Group("")
	eJwt.Use(middleware.JWT([]byte(constant.SECRET_JWT)))
	eJwt.PUT("/admin/:id", controller.UpdateUserAdminController)
	eJwt.DELETE("/admin/:id", controller.DeleteUserAdminController)
	eJwt.GET("/admin", controller.GetUsersAdminController)
	eJwt.GET("/admin/:id", controller.GetUserByidAdminController)
	eJwt.PUT("/user", controller.UpdateUserController)
	eJwt.DELETE("/user", controller.DeleteUserController)
	eJwt.GET("/user", controller.GetUserController)
	eJwt.GET("/image", controller.GetImageHandler)

	bookmark := eJwt.Group("/bookmark")

	NewThreadControllers(eJwt)
	NewBookmarkedContoller(bookmark)

	NewCommentControllers(eJwt)

}

func NewThreadControllers(e *echo.Group) {
	e.GET("/admin/threads", controller.GetThreadController)
	e.GET("/threads/:id", controller.GetThreadsIDController)
	e.POST("/threads", controller.CreateThreadsController)
	e.DELETE("/admin/threads/:id", controller.DeleteThreadsControllerAdmin)
	e.DELETE("/threads/:id", controller.DeleteThreadsControllerAdmin)
	e.PUT("/admin/threads/:id", controller.UpdateThreadsControllerAdmin)
	e.PUT("/threads/:id", controller.UpdateThreadsControllerAdmin)
}

func NewBookmarkedContoller(e *echo.Group) {
	e.GET("", controller.GetSaveThreadController)
	e.POST("", controller.CreateSaveThreadsController)
	e.DELETE("/:id", controller.DeleteSaveThreadsController)
}
func NewCommentControllers(e *echo.Group) {
	e.POST("/comment", controller.CreateCommentController)
	e.DELETE("/comment/:id", controller.DeleteCommentsControllerUser)
	e.PUT("/comment/:id", controller.UpdateCommentsControllerUser)
}
