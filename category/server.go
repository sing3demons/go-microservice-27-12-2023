package main

import (
	"context"
	"errors"
	"fmt"
	"net"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/category/repository"
	"github.com/sing3demons/category/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
)

type HttpServer struct {
	*gin.Engine
}

func NewHttpServer() *HttpServer {
	content, _ := os.ReadFile("VERSION.txt")
	if string(content) != "" {
		buildcommit = string(content)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(ctx *gin.Context) {
		// Starting time request
		startTime := time.Now()
		// Processing request
		ctx.Next()
		// End Time request
		endTime := time.Now()
		// execution time
		latencyTime := endTime.Sub(startTime)
		// Request method
		reqMethod := ctx.Request.Method
		// Request route
		reqUri := ctx.Request.RequestURI
		// status code
		statusCode := ctx.Writer.Status()
		// Request IP
		clientIP := ctx.ClientIP()
		log.WithFields(log.Fields{
			"METHOD":    reqMethod,
			"URI":       reqUri,
			"STATUS":    statusCode,
			"LATENCY":   latencyTime,
			"CLIENT_IP": clientIP,
		}).Info("HTTP REQUEST")
		ctx.Next()
	})

	r.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return &HttpServer{r}
}

func (router *HttpServer) StartHttp(serviceName string) {
	addr := os.Getenv("PORT")
	hostName, _ := os.Hostname()
	srv := &http.Server{
		Addr:           ":" + addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	go func() {
		log.WithFields(log.Fields{
			"PORT":        addr,
			"TYPE":        "HTTP",
			"APP_NAME":    serviceName,
			"APP_VERSION": "1.0.0",
			"BUILD_TIME":  buildtime,
			"APP_COMMIT":  buildcommit,
			"HOSTNAME":    hostName,
		}).Info("Starting server...")

		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			// fmt.Printf("server listen err: %v\n", err)
			log.Error("server listen err: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown: ", err)
		os.Exit(1)
	}

	log.Info("server exiting")
}

func (router *HttpServer) StartGRPC(repo repository.CategoryRepository) {
	s := grpc.NewServer()
	port := os.Getenv("GRPC_PORT")
	hostName, _ := os.Hostname()

	go func() {
		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		// evans --proto=./category.proto -p=50051
		service.RegisterCategoryServiceServer(s, service.NewCategoriesService(repo))
		log.WithFields(log.Fields{
			"PORT":        port,
			"TYPE":        "GRPC",
			"APP_NAME":    "category-server_gRPC",
			"APP_VERSION": "1.0.0",
			"BUILD_TIME":  buildtime,
			"APP_COMMIT":  buildcommit,
			"HOSTNAME":    hostName,
		}).Info("gRPC server is running...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}
