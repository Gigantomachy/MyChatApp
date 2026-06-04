package routers

import (
	"MyChatApp/monolithic/internal/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(authController *controllers.AuthController) *gin.Engine {
	r := gin.Default()

	// setup CORS
	config := cors.DefaultConfig()

	config.AllowOrigins = []string{"http://localhost:5173"}

	// explicitly allow credentials
	config.AllowCredentials = true

	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	r.Use(cors.New(config))

	// register routes for login and registration
	r.POST("/register", authController.Register)
	r.POST("/login", authController.Login)

	return r
}
