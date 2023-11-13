package handlerUtil

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserIdFromContext(c *gin.Context) (int, error) {
	Id := c.Value("userId")
	userId, err := strconv.Atoi(fmt.Sprintf("%v", Id))
	return userId, err
}
