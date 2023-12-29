package repository

import (
	"context"
	"os"
	"time"

	"github.com/sing3demons/category/model"
	"github.com/sing3demons/category/utils"
	"github.com/sirupsen/logrus"
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
	host := os.Getenv("CATEGORY_SERVICE_URL")
	for cur.Next(ctx) {
		var category model.Category
		err := cur.Decode(&category)
		if err != nil {
			return nil, err
		}
		category.Href = utils.Href(host, category.Type, category.ID)
		categories = append(categories, category)
	}

	logrus.WithFields(logrus.Fields{
		"RESULT": categories,
		"TYPE":   "FIND_ALL",
		"COUNT":  len(categories),
		"FILTER": filter,
		"OPTS":   opts,
	}).Info("Find all categories")

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

func (c *categoryRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"id":         id,
		"deleteDate": primitive.Null{},
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	update := bson.M{
		"deleteDate": time.Now().In(loc),
	}

	updateResult, err := c.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"ID":     id,
		"RESULT": updateResult,
		"TYPE":   "DELETE",
	}).Info("Delete category")

	return nil
}
