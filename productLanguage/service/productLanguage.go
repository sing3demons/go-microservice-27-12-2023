package service

import (
	"github.com/sing3demons/productLanguage/model"
	"github.com/sing3demons/productLanguage/repository"
)

type ProductLanguageService interface {
	FindOne(id string) (*model.ProductLanguage, error)
	FindAll() ([]model.ProductLanguage, error)
}

type productLanguageService struct {
	repository repository.ProductLanguageRepository
}

func NewProductLanguageService(repository repository.ProductLanguageRepository) ProductLanguageService {
	return &productLanguageService{repository}
}

func (s *productLanguageService) FindOne(id string) (*model.ProductLanguage, error) {
	return s.repository.FindOne(id)
}

func (s *productLanguageService) FindAll() ([]model.ProductLanguage, error) {
	doc := map[string]interface{}{
		"deleteDate": nil,
		"limit":      20,
	}
	return s.repository.FindAll(doc)
}
