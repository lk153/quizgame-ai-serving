package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/lk153/quizgame-ai-serving/internal/adapters/config"
)

func main() {
	originDomain := "*"
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{originDomain}, // Replace with your frontend's URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	//Init App
	c, err := config.New()
	if err != nil {
		log.Fatalf("Get config: %s\n", err)
	}

	db, err := initializeDB(ctx, c.DB)
	if err != nil {
		log.Fatalf("initializeDB: %s\n", err)
	}

	if c.App.IsCacheOn != config.CACHE_ON {
		c.Redis = nil
	}
	_ = initializeHandlers(ctx, r.Group("/v1"), db, c.Redis)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", c.App.Port),
		Handler: r.Handler(),
	}

	/* Initializing the server in a goroutine so that
	it won't block the graceful shutdown handling below */
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}

		log.Println("App has started on ", c.App.Port)
	}()

	<-ctx.Done()
	stop()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
