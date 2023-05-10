package routes

import (
	"go-simple-products-api/controllers"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.POST("/products", controllers.CreateProduct)
	e.GET("/products", controllers.GetAll)
}
