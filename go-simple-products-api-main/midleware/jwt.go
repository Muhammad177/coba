package midleware

import (
	"Capstone/constant"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func CreateToken(userID int, name string, role string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["name"] = name
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(constant.SECRET_JWT))
}
func ClaimsId(c echo.Context) (float64, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	id := claims["user_id"].(float64)
	return id, nil
}
func ClaimsRole(c echo.Context) (string, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	role := claims["role"].(string)
	return role, nil
}

type TokenStore struct {
	InvalidTokens map[string]time.Time
}

// tokenStore adalah variabel global untuk penyimpanan token
var tokenStore TokenStore

// InitTokenStore inisialisasi penyimpanan token
func InitTokenStore() {
	tokenStore = TokenStore{
		InvalidTokens: make(map[string]time.Time),
	}
}

// DestroyToken menghancurkan token dengan menambahkannya ke daftar token yang tidak valid
func DestroyToken(token string) error {
	// Memastikan token tidak kosong
	if token == "" {
		return errors.New("empty token")
	}

	// Menambahkan token ke daftar token yang tidak valid dengan waktu saat ini
	tokenStore.InvalidTokens[token] = time.Now()

	return nil
}

// IsTokenValid memeriksa apakah token masih valid atau tidak
func IsTokenValid(token string) bool {
	// Memastikan token tidak kosong
	if token == "" {
		return false
	}

	// Memeriksa apakah token ada dalam daftar token yang tidak valid
	_, exists := tokenStore.InvalidTokens[token]
	if exists {
		return false
	}

	return true
}
