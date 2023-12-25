package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/productLanguage/service"
)

type ProductLanguageHandler interface {
	FindAll(c *gin.Context)
	FindOne(c *gin.Context)
}

type productLanguageHandler struct {
	svc service.ProductLanguageService
}

func NewProductLanguageHandler(svc service.ProductLanguageService) ProductLanguageHandler {
	return &productLanguageHandler{svc}
}

func (h *productLanguageHandler) FindAll(c *gin.Context) {
	start := time.Now()

	result, err := h.svc.FindAll()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
		return
	}
	c.JSON(http.StatusOK, result)
	fmt.Println("Time :: ", time.Since(start).String())
}

func (h *productLanguageHandler) FindOne(c *gin.Context) {
	start := time.Now()

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "bad request",
		})
		return
	}
	result, err := h.svc.FindOne(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
		return
	}
	c.JSON(http.StatusOK, result)
	fmt.Println("Time :: ", time.Since(start).String())
}
