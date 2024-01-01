package service

import (
	"github.com/sing3demons/category/model"
	"github.com/sing3demons/category/repository"
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

	return categories, nil
}

func (c *categoryService) FindByID(id string) (*model.Category, error) {
	category, err := c.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return category, nil
}


