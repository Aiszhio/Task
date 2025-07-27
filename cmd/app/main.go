// @title           Subscriptions API
// @version         1.0
// @description     REST‑сервис для учёта онлайн‑подписок
// @host            localhost:8080
// @BasePath        /
//
// @schemes         http
package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	_ "github.com/Aiszhio/Task/docs"
	"github.com/Aiszhio/Task/internal/db"
	"github.com/Aiszhio/Task/internal/middleware"
	"github.com/Aiszhio/Task/internal/repository/postgres"
	transport "github.com/Aiszhio/Task/internal/transport/http"
	"github.com/Aiszhio/Task/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	ctx := context.Background()
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	pool, err := db.NewPool("DATABASE_DSN", ctx)
	if err != nil {
		log.Fatal(err)
	}

	pgRepo := postgres.NewRepository(pool, logger)

	repo := usecase.NewSubscriptionUseCase(pgRepo)

	handlers := transport.NewSubscriptionHandler(repo)

	router := gin.Default()

	router.Use(middleware.NewMiddleware(logger))

	router.Group("/")
	{
		router.GET("/subscriptions/:id", handlers.GetByID())
		router.POST("/subscriptions", handlers.Create())
		router.PUT("/subscriptions/:id", handlers.Update())
		router.DELETE("/subscriptions/:id", handlers.DeleteByID())
		router.POST("/subscriptions/list", handlers.ListByPeriod())
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	err = router.Run(":8080")
	log.Println("Server is running on port 8080")
	if err != nil {
		return
	}
}
