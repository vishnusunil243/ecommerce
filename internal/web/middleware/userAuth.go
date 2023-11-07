package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserAuth(c *gin.Context) {
	// s := c.Request.Header.Get("Authorization")
	tokenString, err := c.Cookie("UserAuth")
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userId, err := ValidateToken(tokenString)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("userId", userId)
	c.Next()
}
