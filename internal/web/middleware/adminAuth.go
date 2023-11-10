package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuth(c *gin.Context) {
	tokenString, err := c.Cookie("AdminAuth")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	adminId, err := ValidateToken(tokenString)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("adminId", adminId)
	c.Next()
}
