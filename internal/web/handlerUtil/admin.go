package handlerUtil

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAdminIdFromContext(c *gin.Context) (int, error) {
	Id := c.Value("adminId")
	adminId, err := strconv.Atoi(fmt.Sprintf("%v", Id))
	return adminId, err
}
