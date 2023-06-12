package midleware

import (
	"Capstone/database"
	"Capstone/models"
	"fmt"

	"github.com/labstack/echo/v4"
)

func BasicAuthDb(username, password string, c echo.Context) (bool, error) {
	var user models.User
	err := database.DB.Where("email = ? AND password = ?", username, password).First(&user).Error
	if err != nil {
		return false, fmt.Errorf("invalid username or password")
	}
	return true, nil
}
