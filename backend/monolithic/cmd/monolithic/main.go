package main

import (
	"context"
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

	session, err := db.InitSession(hosts, keyspace)
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	defer session.Close()

	userRepo := db.NewUserRepository(session)
	userService := services.NewUserService(userRepo)
	authController := controllers.NewAuthController(userService)

	hub := ws.NewHub()

	// optional Redis pub/sub fanout, skip if REDIS_ADDR is unset (single-node
	// run or local tests) - the hub just delivers locally
	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		broker := ws.NewBroker(redisAddr)
		ctx := context.Background()

		if err := broker.Ping(ctx); err != nil {
			log.Printf("WARNING: redis at %s unreachable, ws events will only reach this pod: %v", redisAddr, err)
		}

		hub.AttachBroker(broker)
		broker.Subscribe(ctx, hub)
		defer broker.Close()
	}

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
