package repository

import (
	"context"
	"time"

	"github.com/sing3demons/productLanguage/model"
	"github.com/sing3demons/productLanguage/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductLanguageRepository interface {
	FindOne(id string) (*model.ProductLanguage, error)
	FindAll(doc map[string]any) ([]model.ProductLanguage, error)
}

type productLanguageRepository struct {
	db *mongo.Collection
}

func NewProductLanguageRepository(db *mongo.Collection) ProductLanguageRepository {
	return &productLanguageRepository{db: db}
}

func (r *productLanguageRepository) FindAll(doc map[string]any) ([]model.ProductLanguage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := []model.ProductLanguage{}
	filter, opts := r.QueryBuilder(doc)

	cur, err := r.db.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var elem model.ProductLanguage
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		elem.Href = utils.HostName(elem.Type, elem.ID)

		result = append(result, elem)
	}
	return result, nil
}

func (r *productLanguageRepository) FindOne(id string) (*model.ProductLanguage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result model.ProductLanguage

	filter := bson.M{
		"deleteDate": nil,
		"id":         id,
	}

	if err := r.db.FindOne(ctx, filter).Decode(&result); err != nil {
		return nil, err
	}

	result.Href = utils.HostName(result.Type, result.ID)

	return &result, nil
}

func (r *productLanguageRepository) QueryBuilder(doc map[string]any) (any, *options.FindOptions) {
	var filter bson.D
	opts := options.FindOptions{}
	filter = append(filter, bson.E{Key: "deleteDate", Value: nil})

	for k, v := range doc {
		if k == "limit" {
			limit := v.(int)
			opts.SetLimit(int64(limit))
		}

		if k == "name" {
			filter = append(filter, bson.E{Key: k, Value: v})
		}

	}

	return filter, &opts
}
