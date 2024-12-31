package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lits-06/sell_technology/pkg/utils"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Authorization header missing"})
			c.Abort()
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }

		token, err := utils.VerifyJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid token claims"})
			c.Abort()
			return
		}

		email, emailOk := claims["email"].(string)
		if !emailOk {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid token data"})
			c.Abort()
			return
		}

		c.Set("email", email)
		c.Next()
	}
}