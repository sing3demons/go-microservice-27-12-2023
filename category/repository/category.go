package repository

import "go.mongodb.org/mongo-driver/mongo"

type CategoryRepository interface{}

type categoryRepository struct {
	db *mongo.Collection
}

func NewCategoryRepository(db *mongo.Collection) CategoryRepository {
	return &categoryRepository{db}
}
