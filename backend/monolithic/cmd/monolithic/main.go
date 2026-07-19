package main

import (
	"log"
	"os"
	"strings"

	"MyChatApp/monolithic/internal/controllers"
	"MyChatApp/monolithic/internal/repositories/db"
	"MyChatApp/monolithic/internal/routers"
	"MyChatApp/monolithic/internal/services"
	"MyChatApp/monolithic/internal/ws"
)

func main() {
	hosts := []string{"127.0.0.1"}
	if env := os.Getenv("CASSANDRA_HOSTS"); env != "" {
		hosts = strings.Split(env, ",")
	}

	keyspace := os.Getenv("CASSANDRA_KEYSPACE")
	if keyspace == "" {
		keyspace = "my_chat_app"
	}

	// initialize cassandra DB
	session, err := db.InitSession(hosts, keyspace)
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	defer session.Close()

	// initialize and inject auth controller
	userRepo := db.NewUserRepository(session)
	userService := services.NewUserService(userRepo)
	authController := controllers.NewAuthController(userService)

	hub := ws.NewHub()

	friendRepo := db.NewFriendshipRepository(session)
	friendService := services.NewFriendService(friendRepo, userRepo)
	friendController := controllers.NewFriendController(friendService, hub)

	userController := controllers.NewUserController(userService)

	channelsRepo := db.NewChannelsRepository(session)
	channelsService := services.NewChannelService(channelsRepo, userRepo, friendRepo)
	channelsController := controllers.NewChannelsController(channelsService, hub)

	messagesRepo := db.NewMessagesRepository(session)
	messagesService := services.NewMessageService(messagesRepo, channelsRepo, userRepo)
	messagesController := controllers.NewMessagesController(messagesService, hub)

	r := routers.NewRouter(routers.RouterDependencies{
		Hub:               hub,
		AuthController:    authController,
		FriendController:  friendController,
		UserController:    userController,
		ChannelController: channelsController,
		MessageController: messagesController,
	})

	engine := r.Setup()
	engine.Run(":8080")
}
