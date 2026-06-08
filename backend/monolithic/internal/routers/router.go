package routers

import (
	"MyChatApp/monolithic/internal/controllers"
	"MyChatApp/monolithic/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authController *controllers.AuthController
}

func NewRouter(authcontroller *controllers.AuthController) *Router {
	return &Router{
		authController: authcontroller,
	}
}

func (r *Router) Setup() *gin.Engine {
	engine := gin.Default()

	// setup CORS
	config := cors.DefaultConfig()

	config.AllowOrigins = []string{"http://localhost:5173"}

	// explicitly allow credentials
	config.AllowCredentials = true

	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	engine.Use(cors.New(config))

	// register routes for login and registration
	engine.POST("/register", r.authController.Register)
	engine.POST("/login", r.authController.Login)

	engine.Use(middleware.JWTAuthMiddleware())

	return engine
}
