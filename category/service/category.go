package service

import (
	"os"

	"github.com/sing3demons/category/model"
	"github.com/sing3demons/category/repository"
	"github.com/sing3demons/category/utils"
)

type CategoryService interface {
	FindAll() ([]model.Category, error)
	FindByID(id string) (*model.Category, error)
}

type categoryService struct {
	repository repository.CategoryRepository
}

func NewCategoryService(repository repository.CategoryRepository) CategoryService {
	return &categoryService{repository}
}

func (c *categoryService) FindAll() ([]model.Category, error) {
	categories, err := c.repository.FindAll()
	if err != nil {
		return nil, err
	}

	result := []model.Category{}

	host := os.Getenv("CATEGORY_SERVICE_URL")
	for _, category := range categories {
		category.Href = utils.Href(host, category.Type, category.ID)
		result = append(result, category)
	}
	return result, nil
}

func (c *categoryService) FindByID(id string) (*model.Category, error) {
	category, err := c.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	host := os.Getenv("CATEGORY_SERVICE_URL")
	category.Href = utils.Href(host, category.Type, category.ID)

	return c.repository.FindByID(id)
}
