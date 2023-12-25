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
	"github.com/sing3demons/productLanguage/handler"
	"github.com/sing3demons/productLanguage/repository"
	"github.com/sing3demons/productLanguage/service"
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
	col, err := ConnectMonoDB()
	if err != nil {
		panic(err)
	}

	_ = col
	port := os.Getenv("PORT")
	repositoryProductLanguage := repository.NewProductLanguageRepository(col)
	serviceProductLanguage := service.NewProductLanguageService(repositoryProductLanguage)
	handlerLanguage := handler.NewProductLanguageHandler(serviceProductLanguage)

	r := gin.Default()

	r.GET("/productLanguage", handlerLanguage.FindAll)
	r.GET("/productLanguage/:id", handlerLanguage.FindOne)

	runServer(port, "productLanguage", r)
}

func runServer(addr, serviceName string, router http.Handler) {
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
