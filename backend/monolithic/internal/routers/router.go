package routers

import (
	"MyChatApp/monolithic/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authController *controllers.AuthController) *gin.Engine {
	r := gin.Default()

	r.POST("/register", authController.Register)
	r.POST("/login", authController.Login)

	return r
}