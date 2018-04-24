package middleware

import (
	"git.containerum.net/ch/solutions/pkg/server"
	"github.com/gin-gonic/gin"
)

const (
	//SolutionsServices is key for services
	SolutionsServices = "s-service"
)

// RegisterServices adds services to context
func RegisterServices(svc *server.SolutionsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(SolutionsServices, *svc)
	}
}
