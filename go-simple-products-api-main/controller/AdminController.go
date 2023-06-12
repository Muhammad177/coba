package controller

import (
	"Capstone/database"
	"Capstone/midleware"
	"Capstone/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo/v4"
)

func CreatePhoto(c echo.Context) (string, error) {
	// Check if a file photo is present in the request
	_, err := c.FormFile("photo")
	if err != nil {
		return "", nil
	}

	// Menerima file foto dari permintaan
	file, err := c.FormFile("photo")
	if err != nil {
		return "", echo.NewHTTPError(http.StatusBadRequest, "Failed to upload photo")
	}

	// Generate nama unik untuk file foto
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Buka file foto yang diunggah
	src, err := file.Open()
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, "Failed to open photo")
	}
	defer src.Close()

	// Simpan file foto di direktori lokal
	dstPath := "uploads/" + filename
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, "Failed to save photo")
	}
	defer dst.Close()

	// Salin isi file foto yang diunggah ke file tujuan
	if _, err = io.Copy(dst, src); err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, "Failed to save photo")
	}

	// Mengembalikan path file foto
	return dstPath, nil
}

func CreateUserController(c echo.Context) error {
	// Bind data pengguna dari permintaan
	user := models.User{}
	err := c.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := database.DB.Where("email = ?", user.Email).First(&user).Error; err == nil {
		// Email sudah ada, kembalikan respons error
		return echo.NewHTTPError(http.StatusBadRequest, "Email already exists")
	}
	// Set nilai default untuk role
	user.Role = "User"

	// Simpan foto pengguna
	photoPath, err := CreatePhoto(c)
	if err != nil {
		return err
	}

	// Set path file foto pengguna
	user.Photo = photoPath

	// Simpan pengguna ke database
	err = database.DB.Save(&user).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save user")
	}

	// Mengembalikan respons JSON
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Berhasil membuat pengguna baru",
		"user":    user,
	})
}
func UpdateUserAdminController(c echo.Context) error {
	role, err := midleware.ClaimsRole(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, "Only admin can access")
	}

	id := c.Param("id")

	var users models.User
	if err := database.DB.Where("id = ?", id).First(&users).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	previousEmail := users.Email

	if err := c.Bind(&users); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if previousEmail != users.Email {
		var existingUser models.User
		if err := database.DB.Where("email = ?", users.Email).First(&existingUser).Error; err == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Email already exists")
		}
	}

	// Check if new photo is uploaded
	_, err = c.FormFile("photo")
	if err == nil {
		// New photo is uploaded, execute CreatePhoto function
		photoPath, err := CreatePhoto(c)
		if err != nil {
			return err
		}

		// Delete previous photo
		if users.Photo != "" {
			if err := os.Remove(users.Photo); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user photo")
			}
		}

		users.Photo = photoPath
	} else if err == http.ErrMissingFile {
		// No new photo provided, check if existing photo needs to be deleted
		if users.Photo != "" {
			// Delete previous photo from database and local directory
			if err := os.Remove(users.Photo); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user photo")
			}

			users.Photo = ""
		}
	}

	// Update the user in the database
	if err := database.DB.Model(&users).Updates(map[string]interface{}{
		"photo": users.Photo,
	}).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User updated successfully",
		"user":    users,
	})
}

func GetUserByidAdminController(c echo.Context) error {
	role, err := midleware.ClaimsRole(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, "Only admin can access")
	}

	id := c.Param("id")
	var users models.User
	if err := database.DB.Where("id = ?", id).First(&users).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get User by id",
		"user":    users,
	})
}
func GetUsersAdminController(c echo.Context) error {
	role, err := midleware.ClaimsRole(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, "Only admin can access")
	}

	var users []models.User
	err = database.DB.Find(&users).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve users from the database")
	}
	allUsers := make([]models.AllUser, len(users))
	for i, user := range users {
		allUsers[i] = models.ConvertUserToAllUser(&user)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success: Retrieved all users",
		"users":   allUsers,
	})
}
func DeleteUserAdminController(c echo.Context) error {
	role, err := midleware.ClaimsRole(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, "Only admin can access")
	}

	id := c.Param("id")
	var users models.User

	if err := database.DB.Where("id = ?", id).First(&users).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Menghapus file foto terkait jika ada
	if users.Photo != "" {
		if err := os.Remove(users.Photo); err != nil {
			// Jika gagal menghapus file, Anda dapat menangani kesalahan di sini
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user photo")
		}
	}

	if err := database.DB.Delete(&users).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success delete user by ID",
		"user":    users,
	})
}

func LoginAdminController(c echo.Context) error {
	admin := models.AdminResponse{ID: 1, Name: "Wahyu", Email: "admin@gmail.com", Password: "admin123"}
	if err := c.Bind(&admin); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}
	var admins = []models.AdminResponse{
		{ID: 1, Name: "Wahyu", Email: "admin@gmail.com", Password: "admin123"},
	}
	for _, a := range admins {
		if a.Email == admin.Email && a.Password == admin.Password {
			token, err := midleware.CreateToken(int(a.ID), a.Name, "admin")
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"message": "Failed Login",
					"error":   err.Error(),
				})
			}

			adminResponse := models.UserResponse{admin.ID, admin.Name, admin.Email, token}
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": "Login Admin Sukses",
				"Admin":   adminResponse,
			})
		}
	}
	return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid username or password"})
}
