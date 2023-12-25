package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("ZONE") == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		if err := godotenv.Load(".env.dev"); err != nil {
			panic(err)
		}
	}
}

func main() {
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live")

	connect, err := NewMongo()
	if err != nil {
		panic(err)
	}

	// InitSeed(connect)
	col := connect.Database("category").Collection("category")
	_ = col

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	runServer("category", r)
}

func runServer(serviceName string, router http.Handler) {
	addr := os.Getenv("PORT")
	srv := &http.Server{
		Addr:    ":" + addr,
		Handler: router,
	}

	go func() {
		fmt.Printf("[%s] http listen: %s\n", serviceName, srv.Addr)

		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server listen err: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown: ", err)
	}

	fmt.Println("server exited")
}