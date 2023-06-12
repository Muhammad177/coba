package controller

import (
	"fmt"
	"net/http"

	"Capstone/midleware"
	"strconv"

	"Capstone/database"
	// "Capstone/models"

	"github.com/labstack/echo/v4"
)

func GetSaveThreadController(c echo.Context) error {
	id, err := midleware.ClaimsId(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	svThread, err := database.GetSaveThreads(c.Request().Context(), int(id))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success getting Bookmark",
		"data":    svThread,
	})

}

type BookmarkThreadRequest struct {
	ID int `json:"thread_id"`
}

func CreateSaveThreadsController(c echo.Context) error {
	svThread := BookmarkThreadRequest{}
	if err := c.Bind(&svThread); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Print(svThread)

	id, err := midleware.ClaimsId(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	newSvThread, err := database.CreateSaveThreads(c.Request().Context(), int(id), svThread.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success creating thread",
		"data":    newSvThread,
	})
}

func DeleteSaveThreadsController(c echo.Context) error {
	userID, err := midleware.ClaimsId(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	threadID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = database.DeleteSaveThreads(c.Request().Context(), int(userID), threadID)
	if err != nil {
		if err == database.ErrInvalidID {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success unbookmark thread",
	})
}
