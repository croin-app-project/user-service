package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/croin-app-project/user-service/internal/domain"

	http_response "github.com/croin-app-project/package/pkg/utils/http-response"
	_config "github.com/croin-app-project/user-service/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Authorized(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader != "" {

		tokenString := strings.Split(authHeader, " ")[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			config := _config.ReadConfiguration()
			return []byte(config.App.Jwt.Key), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok {

			// Check expired
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				errCode, errOjb := http_response.HandleException(http_response.AUTHORIZATION_KEY_EXPIRED, err)
				return c.Status(errCode).JSON(errOjb)
			}

			// c.Set("username", claims["sub"])
		} else {
			errCode, errOjb := http_response.HandleException(http_response.INVALID_AUTHORIZATION_KEY, err)
			return c.Status(errCode).JSON(errOjb)
		}
	} else {
		err := errors.New("token is required")
		errCode, errOjb := http_response.HandleException(http_response.AUTHORIZATION_KEY_INACTIVE, err)
		return c.Status(errCode).JSON(errOjb)
	}
	return c.Next()
}

func BuildToken(user domain.User) (string, error) {
	config := _config.ReadConfiguration()
	secretKey := []byte(config.App.Jwt.Key)
	expMinutes := config.App.Jwt.Exp
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.Username,
		"username": user.Username,
		"name":     user.Username,
		"exp":      time.Now().Add(time.Minute * time.Duration(expMinutes)).Unix(),
	})
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}
