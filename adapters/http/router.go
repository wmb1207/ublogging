package http

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wmb1207/ublogging/adapters/http/handler"
	"github.com/wmb1207/ublogging/config"
)

type WithRouter func(router *gin.RouterGroup)

func Router(cfg *config.Config, env string, routers ...WithRouter) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(
		cors.Config{
			AllowOrigins: cfg.AllowedOrigins,
			AllowMethods: []string{
				http.MethodGet,
				http.MethodPut,
				http.MethodPost,
				http.MethodOptions,
				http.MethodDelete,
			},
			AllowHeaders:     cfg.AllowedHeaders,
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		},
	))

	r.Handle(http.MethodGet, "/health", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]bool{"health": true})
	})

	api := r.Group("/api")
	for _, router := range routers {
		router(api)
	}

	return r
}

func WithUserRouter(userHandler *handler.UserHandler, middlewares ...gin.HandlerFunc) WithRouter {
	return func(router *gin.RouterGroup) {
		users := router.Group("/users")
		users.POST("", userHandler.New)
		users.GET("/feed", append(middlewares, userHandler.Feed)...)
		users.POST("/:user_uuid/follow", append(middlewares, userHandler.Follow)...)
	}
}

func WithPostRouter(postHandler *handler.PostHandler, middlewares ...gin.HandlerFunc) WithRouter {
	return func(router *gin.RouterGroup) {
		posts := router.Group("/posts")
		posts.POST("", append(middlewares, postHandler.New)...)
		posts.GET("/:post_uuid", append(middlewares, postHandler.Post)...)
		posts.POST("/:post_uuid", append(middlewares, postHandler.Comment)...)
	}

}
