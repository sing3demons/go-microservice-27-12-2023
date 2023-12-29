package main

import (
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sing3demons/category/handler"
	"github.com/sing3demons/category/repository"
	"github.com/sing3demons/category/service"
	"github.com/sing3demons/category/store"
	log "github.com/sirupsen/logrus"
)

func runCMD(cmd string, shell bool) []byte {
	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
			panic("some error found")
		}
		return out
	}
	out, err := exec.Command(cmd).Output()
	if err != nil {
		log.Fatal(err)
	}

	return out
}

func init() {

	if os.Getenv("ZONE") == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		if err := godotenv.Load(".env.dev"); err != nil {
			panic(err)
		}
	}

	// setup logrus
	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = log.InfoLevel
	}

	log.SetLevel(logLevel)
	log.SetFormatter(&log.JSONFormatter{})

}

func main() {
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer os.Remove("/tmp/live")

	col := store.NewStore().Client.Database("category").Collection("category")
	repo := repository.NewCategoryRepository(col)
	svc := service.NewCategoryService(repo)
	controller := handler.NewCategoryHandler(svc)

	r := NewHttpServer()

	r.GET("/category", controller.FindAll)
	r.GET("/category/:id", controller.FindByID)

	r.StartHttp("category")
}
