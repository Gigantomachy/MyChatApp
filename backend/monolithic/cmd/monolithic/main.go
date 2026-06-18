package main

import (
	"log"
	"os"
	"strings"

	"MyChatApp/monolithic/internal/controllers"
	"MyChatApp/monolithic/internal/repositories/db"
	"MyChatApp/monolithic/internal/routers"
	"MyChatApp/monolithic/internal/services"
)

func main() {
	hosts := []string{"127.0.0.1"}
	if env := os.Getenv("CASSANDRA_HOSTS"); env != "" {
		hosts = strings.Split(env, ",")
	}

	// initialize cassandra DB
	session, err := db.InitSession(hosts, "my_chat_app")
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	defer session.Close()

	// initialize and inject auth controller
	userRepo := db.NewUserRepository(session)
	userService := services.NewUserService(userRepo)
	authController := controllers.NewAuthController(userService)

	friendRepo := db.NewFriendshipRepository(session)
	friendService := services.NewFriendService(friendRepo, userRepo)
	friendController := controllers.NewFriendController(friendService)

	userController := controllers.NewUserController(userService)

	r := routers.NewRouter(routers.RouterDependencies{
		AuthController:   authController,
		FriendController: friendController,
		UserController:   userController,
	})

	engine := r.Setup()
	engine.Run(":8080")
}
