package repository

import (
	"context"
	"time"

	"github.com/sing3demons/category/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryRepository interface {
	FindAll() ([]model.Category, error)
	FindByID(id string) (*model.Category, error)
}

type categoryRepository struct {
	db *mongo.Collection
}

func NewCategoryRepository(db *mongo.Collection) CategoryRepository {
	return &categoryRepository{db}
}

func (c *categoryRepository) FindAll() ([]model.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	categories := []model.Category{}

	filter := bson.M{
		"deleteDate": primitive.Null{},
	}
	opts := options.FindOptions{}
	opts.SetProjection(bson.M{
		"_id":      0,
		"products": 0,
	})
	cur, err := c.db.Find(ctx, filter, &opts)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var category model.Category
		err := cur.Decode(&category)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (c *categoryRepository) FindByID(id string) (*model.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var category model.Category

	filter := bson.M{
		"id":         id,
		"deleteDate": primitive.Null{},
	}

	err := c.db.FindOne(ctx, filter).Decode(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}
