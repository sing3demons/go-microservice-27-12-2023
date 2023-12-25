package service

import "github.com/sing3demons/category/repository"

type CategoryService interface{}

type categoryService struct {
	repository repository.CategoryRepository
}

func NewCategoryService(repository repository.CategoryRepository) CategoryService {
	return &categoryService{repository}
}
