package server

import (
	"github.com/gin-gonic/gin"
)

func tenantUUIDFromContext(c *gin.Context) string {
	return c.Value("tenantUUID").(string)
}
