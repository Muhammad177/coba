package controller

import (
	"net/http"
	"strconv"

	"Capstone/database"
	"Capstone/midleware"
	"Capstone/models"

	"github.com/labstack/echo/v4"
)

func CreateFollowController(c echo.Context) error {
	Follow := models.Follow{}
	c.Bind(&Follow)
	id, err := midleware.ClaimsId(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	Follow.UserID = int(id)
	newFollow, err := database.CreateFollow(c.Request().Context(), Follow)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success creating Follow",
		"data":    newFollow,
	})
}

func DeleteFollowsControllerUser(c echo.Context) error {
	id, err := midleware.ClaimsId(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var users models.User
	if err := database.DB.Where("id = ?", id).First(&users).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = database.DeleteFollows(c.Request().Context(), Id)
	if err != nil {
		if err == database.ErrInvalidID {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success deleting Follow data",
	})
}
