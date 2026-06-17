package routers

import (
	"MyChatApp/monolithic/internal/controllers"
	"MyChatApp/monolithic/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// use a dependencies struct to avoid having to pass a million constructor arguments
type RouterDependencies struct {
	AuthController   *controllers.AuthController
	FriendController *controllers.FriendController

	// add channel controller later
}

type Router struct {
	RouterDependencies
}

func NewRouter(deps RouterDependencies) *Router {
	return &Router{deps}
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

	api := engine.Group("/api") // normal assignment

	{ // braces for readability

		// register routes for login and registration
		api.POST("/register", r.AuthController.Register)
		api.POST("/login", r.AuthController.Login)

		api.Use(middleware.JWTAuthMiddleware())

		{
			api.GET("/friends", r.FriendController.GetFriends)
			api.DELETE("/friends/:id", r.FriendController.RemoveFriend)

			api.GET("/friend-requests", r.FriendController.GetFriendRequests)
			api.GET("/friend-requests/:id", r.FriendController.GetFriendRequestByID)
			api.POST("/friend-requests/:id", r.FriendController.SendFriendRequest)
			api.PUT("/friend-requests/:id", r.FriendController.AcceptFriendRequest)
			api.DELETE("/friend-requests/outgoing/:id", r.FriendController.CancelFriendRequest)
			api.DELETE("/friend-requests/incoming/:id", r.FriendController.RejectFriendRequest)
		}
	}

	return engine
}
