package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/product/service"
)

type ProductHandler interface {
	FindAll(c *gin.Context)
	FindOne(c *gin.Context)
	FindProductLanguage(c *gin.Context)
	FindProductLanguages(c *gin.Context)
}

type productHandler struct {
	svc service.ProductService
}

func NewProductHandler(svc service.ProductService) ProductHandler {
	return &productHandler{svc}
}

func (h *productHandler) FindAll(c *gin.Context) {
	start := time.Now()
	name := c.Query("name")

	query := map[string]any{}
	if name != "" {
		query["name"] = name
	}

	sort := c.DefaultQuery("sort", "_id")
	query["sort"] = sort

	order := c.DefaultQuery("order", "desc")
	query["order"] = order

	limit := c.DefaultQuery("limit", "100")
	size, err := strconv.Atoi(limit)
	if err != nil {
		size = 10
	}
	query["limit"] = size

	skip := c.DefaultQuery("skip", "0")
	skipInt, err := strconv.Atoi(skip)
	if err != nil {
		skipInt = 0
	}
	query["skip"] = skipInt

	project := c.Query("project")
	if project != "" {
		query["project"] = project
	}

	expand := c.Query("expand")
	if expand != "" {
		expands := strings.Split(expand, ",")
		if len(expands) > 0 {
			for _, v := range expands {
				if strings.Contains(v, "product.productLanguage") {
					query["language"] = true
				}
			}
		}
	}

	products, total, err := h.svc.FindAll(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"duration": time.Since(start).String(),
		"products": products,
		"total":    total,
	})
}

func (h *productHandler) FindOne(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "bad request",
		})
		return
	}
	product, err := h.svc.FindOne(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
		return
	}
	c.JSON(http.StatusOK, product)
	fmt.Println("Time :: ", time.Since(start))
}

func (h *productHandler) FindProductLanguage(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	languageCode := c.Param("languageCode")
	if id == "" || languageCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "bad request",
		})
		return
	}

	productLanguage, err := h.svc.FindProductLanguage(id, languageCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
		return
	}
	c.JSON(http.StatusOK, productLanguage)
	fmt.Println("Time :: ", time.Since(start))
}

func (h *productHandler) FindProductLanguages(c *gin.Context) {
	start := time.Now()
	languageCode := c.Param("languageCode")
	productLanguage, err := h.svc.FindProductLanguages(languageCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
		return
	}
	c.JSON(http.StatusOK, productLanguage)
	fmt.Println("Time :: ", time.Since(start))
}
