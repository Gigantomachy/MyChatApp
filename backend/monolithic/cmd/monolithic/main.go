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

	session, err := db.InitSession(hosts, "my_chat_app")
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	defer session.Close()

	userRepo := db.NewUserRepository(session)
	userService := services.NewUserService(userRepo)
	authController := controllers.NewAuthController(userService)

	r := routers.SetupRouter(authController)
	r.Run(":8080")
}