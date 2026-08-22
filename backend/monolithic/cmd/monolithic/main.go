package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"MyChatApp/monolithic/internal/controllers"
	"MyChatApp/monolithic/internal/repositories/db"
	"MyChatApp/monolithic/internal/routers"
	"MyChatApp/monolithic/internal/services"
	"MyChatApp/monolithic/internal/ws"
)

func main() {
	hosts := []string{"127.0.0.1"}
	if env := os.Getenv("CASSANDRA_HOSTS"); env != "" {
		parts := strings.Split(env, ",")
		hosts = hosts[:0] // reset to empty, keep capacity
		for _, p := range parts {
			if p = strings.TrimSpace(p); p != "" {
				hosts = append(hosts, p)
			}
		}
	}

	keyspace := os.Getenv("CASSANDRA_KEYSPACE")
	if keyspace == "" {
		keyspace = "my_chat_app"
	}

	db_username := os.Getenv("CASSANDRA_USERNAME")
	db_password := os.Getenv("CASSANDRA_PASSWORD")

	session, err := db.InitSession(hosts, keyspace, db_username, db_password)
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	// defer session.Close()

	userRepo := db.NewUserRepository(session)
	userService := services.NewUserService(userRepo)
	authController := controllers.NewAuthController(userService)

	hub := ws.NewHub()

	// optional Redis pub/sub fanout, skip if REDIS_ADDR is unset (single-node
	// run or local tests) - the hub just delivers locally
	var broker *ws.Broker
	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		broker = ws.NewBroker(redisAddr)
		ctx := context.Background()

		if err := broker.Ping(ctx); err != nil {
			log.Fatalf("redis at %s unreachable: %v", redisAddr, err)
		}

		hub.AttachBroker(broker)
		broker.Subscribe(ctx, hub)
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

	healthController := controllers.NewHealthController(session)

	r := routers.NewRouter(routers.RouterDependencies{
		Hub:               hub,
		AuthController:    authController,
		FriendController:  friendController,
		UserController:    userController,
		ChannelController: channelsController,
		MessageController: messagesController,
		HealthController:  healthController,
	})

	engine := r.Setup()
	srv := &http.Server{Addr: ":8080", Handler: engine}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	session.Close()
	if broker != nil {
		broker.Close()
	}
}
