package middlewares

import (
	"finalProject2/pkg/helpers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenClaims, err := helpers.VerifyToken(c)
		fmt.Println(tokenClaims)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "user is unauthorized",
				"error":   err.Error(),
			})
			return
		}
		c.Set("user", tokenClaims)
		c.Next()

		return

	}
}
