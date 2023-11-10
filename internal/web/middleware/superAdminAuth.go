package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuperAdminAuth(c *gin.Context) {
	// s := c.Request.Header.Get("Authorization")
	tokenString, err := c.Cookie("SuperAdminAuth")
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	superId, err := ValidateToken(tokenString)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("superId", superId)
	c.Next()
}
