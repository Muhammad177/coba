package controllers

import (
	"go-simple-products-api/models"
	"go-simple-products-api/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response[T any] struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func CreateProduct(c echo.Context) error {
	var product models.Product

	if err := c.Bind(&product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	createdProduct, err := services.CreateProduct(product)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create a product")
	}

	return c.JSON(http.StatusCreated, Response[models.Product]{
		Message: "product created",
		Data:    createdProduct,
	})
}

func GetAll(c echo.Context) error {
	products := services.GetProducts()

	return c.JSON(http.StatusOK, Response[[]models.Product]{
		Message: "all products",
		Data:    products,
	})
}
