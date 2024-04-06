package routes

import (
	"markie-backend/controllers"
	"markie-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SpaceRoutes(rg *gin.RouterGroup) {
	space := rg.Group("/space")
	space.GET("/all", middlewares.VerifyMiddleware(), controllers.GetUserSpaces)
}
