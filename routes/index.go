package routes

import (
	"github.com/gin-gonic/gin"
)

func GetRoutes(r *gin.Engine) {
	route := r.Group("/")
	api := r.Group("/api")
	UserRoutes(api)
	SpaceRoutes(api)
	PromptRoutes(api)
	route.GET("health", healthcheck)
}

func healthcheck(c *gin.Context) {
	c.JSON(200, gin.H{"message": "API has good health"})
}
