package routes

import (
	"markie-backend/controllers"
	"markie-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func PromptRoutes(rg *gin.RouterGroup) {
	prompt := rg.Group("/prompt")
	prompt.POST("/", middlewares.VerifyMiddleware(), controllers.ReceivePrompt)
}
