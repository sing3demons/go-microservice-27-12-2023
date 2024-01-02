package service

import (
	"context"

	"github.com/sing3demons/product/model"
)

type CategoryService interface {
	GetCategory(id string) (*model.Category, error)
}

type categoryService struct {
	categoryClient CategoryServiceClient
}

func NewCategoryService(categoryClient CategoryServiceClient) CategoryService {
	return &categoryService{categoryClient}
}

func (s *categoryService) GetCategory(id string) (*model.Category, error) {
	req := CategoryRequest{
		Id: id,
	}

	res, err := s.categoryClient.GetCategory(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	result := &model.Category{
		Name:       res.Name,
		Type:       res.Type,
		ID:         res.Id,
		Href:       res.Href,
		Version:    res.Version,
		LastUpdate: res.LastUpdate,
	}

	return result, nil
}
