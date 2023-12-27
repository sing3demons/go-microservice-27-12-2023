package handler

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/category/model"
	"github.com/sing3demons/category/service"
	"github.com/sing3demons/category/utils"
)

type CategoryHandler interface {
	FindAll(c *gin.Context)
	FindByID(c *gin.Context)
}

type categoryHandler struct {
	svc service.CategoryService
}

func NewCategoryHandler(svc service.CategoryService) CategoryHandler {
	return &categoryHandler{svc}
}

func (h *categoryHandler) FindAll(c *gin.Context) {
	categories, err := h.svc.FindAll()
	if err != nil {
		c.JSON(404, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(200, categories)
}

func (h *categoryHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	category, err := h.svc.FindByID(id)
	if err != nil {
		c.JSON(404, gin.H{
			"message": err.Error(),
		})
	}

	limit := c.DefaultQuery("limit", "1000")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 100
	}

	skip := c.DefaultQuery("skip", "0")
	skipInt, err := strconv.Atoi(skip)
	if err != nil {
		skipInt = 0
	}

	if len(category.Products) > 0 {
		host := os.Getenv("PRODUCT_SERVICE_URL")
		products := []model.Products{}
		for i, product := range category.Products {
			products = append(products, model.Products{
				Type: "products",
				ID:   product.ID,
				Name: product.Name,
				Href: utils.Href(host, "products", product.ID),
			})

			if i == limitInt {
				break
			}
		}
		category.Products = paginate(products, skipInt, limitInt)
	}

	c.JSON(200, category)
}

func paginate(x []model.Products, skip int, size int) []model.Products {
	if skip > len(x) {
		skip = len(x)
	}

	end := skip + size
	if end > len(x) {
		end = len(x)
	}

	return x[skip:end]
}
