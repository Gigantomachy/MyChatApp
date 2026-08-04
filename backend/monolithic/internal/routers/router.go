package routers

import (
	"MyChatApp/monolithic/internal/controllers"
	"MyChatApp/monolithic/internal/middleware"
	"MyChatApp/monolithic/internal/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

// use a dependencies struct to avoid having to pass a million constructor arguments
type RouterDependencies struct {
	Hub               *ws.Hub
	AuthController    *controllers.AuthController
	FriendController  *controllers.FriendController
	UserController    *controllers.UserController
	ChannelController *controllers.ChannelsController
	MessageController *controllers.MessagesController
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

	config.AllowOrigins = strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

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

		api.POST("/logout", r.AuthController.Logout)
		api.GET("/me", r.AuthController.Me)
		api.GET("/ws", ws.ServeWS(r.Hub))

		{
			api.GET("/users", r.UserController.SearchUsers)

			api.GET("/friends", r.FriendController.GetFriends)
			api.DELETE("/friends/:id", r.FriendController.RemoveFriend)

			api.GET("/friend-requests", r.FriendController.GetFriendRequests)
			api.GET("/friend-requests/:id", r.FriendController.GetFriendRequestByID)
			api.POST("/friend-requests/:id", r.FriendController.SendFriendRequest)
			api.PUT("/friend-requests/:id", r.FriendController.AcceptFriendRequest)
			api.DELETE("/friend-requests/outgoing/:id", r.FriendController.CancelFriendRequest)
			api.DELETE("/friend-requests/incoming/:id", r.FriendController.RejectFriendRequest)

			api.GET("/channels/discover", r.ChannelController.GetAllChannels)
			api.GET("/channels", r.ChannelController.GetChannelsByUser)
			api.GET("/channels/:id", r.ChannelController.GetChannelUsersAndInformation)
			api.POST("/channels", r.ChannelController.CreateChannel)
			api.POST("/channels/:id", r.ChannelController.CreateChannelMembership)
			api.PUT("/channels/:id", r.ChannelController.ModifyChannelMembership)
			api.DELETE("/channels/:id", r.ChannelController.DeleteChannel)
			api.DELETE("/channels/:id/membership", r.ChannelController.DeleteChannelMembership)

			api.GET("/channels/:id/messages", r.MessageController.GetMessages)
			api.POST("/channels/:id/messages", r.MessageController.SendMessage)
		}
	}

	return engine
}
