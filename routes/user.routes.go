package routes

import (
	"markie-backend/controllers"
	"markie-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/auth")
	users.POST("/register", controllers.ResgisterUser)
	users.GET("/users", middlewares.VerifyMiddleware(), controllers.GetUsers)
	users.GET("/verify", middlewares.VerifyMiddleware(), controllers.VerifyUser)
	users.GET("/:id", middlewares.VerifyMiddleware(), controllers.GetUserById)
	users.POST("/login", controllers.Login)
	users.POST("/logout", middlewares.VerifyMiddleware(), controllers.Logout)
	users.GET("/spaces", middlewares.VerifyMiddleware(), controllers.GetUserSpaces)
}
