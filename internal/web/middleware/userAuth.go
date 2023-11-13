package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main.go/internal/web/handlerUtil"
)

func UserAuth(c *gin.Context) {
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
func UserIsBlocked(c *gin.Context) {
	var cr *gorm.DB
	id, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("user is blocked"))
		return
	}
	var isBlocked bool
	err = cr.Raw(`SELECT is_blocked FROM users WHERE id=?`, id).Scan(&isBlocked).Error
	if err != nil {
		c.Abort()
		return
	}
	if isBlocked {
		c.Abort()
		return
	}
	c.Next()
}
