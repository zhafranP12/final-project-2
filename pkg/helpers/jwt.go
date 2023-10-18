package helpers

import (
	"finalProject2/infrastructure/config"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)


func GenerateToken(id int, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"email": email,
	})

	tokenString, err := token.SignedString([]byte(config.GetAppConfig().JWTSecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(c *gin.Context) (interface{}, error) {
	errResponse := errors.New("a token is required to access")
	getToken := c.Request.Header.Get("Authorization")
	bearer := strings.HasPrefix(getToken, "Bearer")

	if !bearer {
		return nil, errResponse
	}

	tokenString := strings.Split(getToken, " ")[1]

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.GetAppConfig().JWTSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}
