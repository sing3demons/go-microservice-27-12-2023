package handler

import "github.com/sing3demons/category/service"

type CategoryHandler interface{}

type categoryHandler struct {
	svc service.CategoryService
}

func NewCategoryHandler(svc service.CategoryService) CategoryHandler {
	return &categoryHandler{svc}
}
