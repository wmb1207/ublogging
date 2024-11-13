package main

import (
	"fmt"
	"strconv"

	mongodb "github.com/wmb1207/ublogging/adapters/database/mongo"
	"github.com/wmb1207/ublogging/adapters/http"
	"github.com/wmb1207/ublogging/adapters/http/handler"
	"github.com/wmb1207/ublogging/config"
	"github.com/wmb1207/ublogging/internal/service"
)

func main() {
	fmt.Println("Starting Ublogging")

	cfg := config.LoadConfig()

	baseRepo, err := mongodb.NewMongoRepository(
		mongodb.WithDBName(cfg.DBName),
		mongodb.WithDBURI(fmt.Sprintf("%s://%s:%s@%s:%s", cfg.DBEngine, cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort)),
	)

	if err != nil {
		panic(err)
	}

	userRepo := mongodb.NewMongoUserRepository(baseRepo)
	postRepo := mongodb.NewMongoPostRepository(baseRepo)

	userService := service.NewUserService(userRepo, postRepo)
	postService := service.NewPostService(postRepo)

	userHandler := handler.NewUserHandler(userService, postService)
	postHandler := handler.NewPostHandler(userService, postService)

	userTokenMiddleware := handler.TokenAuthMiddleware(userService)

	router := http.Router(
		cfg,
		cfg.ENV,
		http.WithUserRouter(userHandler, userTokenMiddleware),
		http.WithPostRouter(postHandler, userTokenMiddleware))

	router.Run("0.0.0.0:" + strconv.Itoa(cfg.Port))
}
