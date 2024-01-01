package service

import (
	context "context"

	"github.com/sing3demons/category/repository"
)

type categoryGrpcService struct {
	repository repository.CategoryRepository
}

func NewCategoriesService(repository repository.CategoryRepository) CategoryServiceServer {
	return &categoryGrpcService{repository: repository}
}

func (s *categoryGrpcService) GetCategory(ctx context.Context, in *CategoryRequest) (*CategoryResponse, error) {
	resp, err := s.repository.FindByID(in.Id)
	if err != nil {
		return nil, err
	}

	return &CategoryResponse{
		Id:         resp.ID,
		Name:       resp.Name,
		Type:       resp.Type,
		Href:       resp.Href,
		LastUpdate: resp.LastUpdate,
		Version:    resp.Version,
	}, nil
}

func (s *categoryGrpcService) mustEmbedUnimplementedCategoryServiceServer() {}
