package controller

import (
	"net/http"
	"strconv"

	"Capstone/database"
	"Capstone/midleware"
	"Capstone/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetThreadController(c echo.Context) error {
	role, err := midleware.ClaimsRole(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, "Only admin can access")
	}
	thread, err := database.GetThreads(c.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	allThreads := make([]models.AllThread, len(thread))
	for i, thread := range thread {
		allThreads[i] = models.ConverThreadToAllThread(&thread)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success: Retrieved all threads",
		"data":    allThreads,
	})
}

func GetThreadsIDController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	thread, err := database.GetThreadsByID(c.Request().Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success getting Thread",
		"data":    models.ConvertThreadToThreadResponse(&thread),
	})
}
func GetThreadControllerByTitle(c echo.Context) error {
	title := c.QueryParam("title")
	thread, err := database.GetThreadByTitle(c.Request().Context(), title)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully retrieved thread by title",
		"thread":  thread,
	})
}

func CreateThreadsController(c echo.Context) error {
	thread := models.Thread{}
	c.Bind(&thread)
	id, err := midleware.ClaimsId(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	thread.UserID = int(id)
	newThread, err := database.CreateThreads(c.Request().Context(), thread)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success creating thread",
		"data":    newThread,
	})
}

func DeleteThreadsControllerAdmin(c echo.Context) error {
	role, err := midleware.ClaimsRole(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, "Only admin can access")
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = database.DeleteThreads(c.Request().Context(), id)
	if err != nil {
		if err == database.ErrInvalidID {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success deleting thread data",
	})
}

func UpdateThreadsControllerAdmin(c echo.Context) error {
	role, err := midleware.ClaimsRole(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, "Only admin can access")
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	thread := models.Thread{}
	c.Bind(&thread)

	updateThread, err := database.UpdateThreads(c.Request().Context(), id, thread)
	if err != nil {
		if err == database.ErrInvalidID {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success updating thread data",
		"data":    updateThread,
	})
}

func DeleteThreadsControllerUser(c echo.Context) error {
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

	err = database.DeleteThreads(c.Request().Context(), Id)
	if err != nil {
		if err == database.ErrInvalidID {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success deleting thread data",
	})
}

func UpdateThreadsControllerUser(c echo.Context) error {
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

	thread := models.Thread{}
	c.Bind(&thread)

	updateThread, err := database.UpdateThreads(c.Request().Context(), Id, thread)
	if err != nil {
		if err == database.ErrInvalidID {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success updating thread data",
		"data":    updateThread,
	})
}
