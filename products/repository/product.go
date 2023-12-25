package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/sing3demons/product/model"
	"github.com/sing3demons/product/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository interface {
	FindAll(doc map[string]any) ([]model.Products, error)
	CountProduct() int64
	FindOne(id string) (*model.Products, error)
	FindProductLanguage(id, languageCode string) (*model.ProductLanguage, error)
	FindProductLanguages(languageCode string) ([]model.ProductLanguage, error)
}

type productRepository struct {
	db *mongo.Collection
}

func NewProducts(db *mongo.Collection) ProductRepository {
	return &productRepository{db}
}

func (p *productRepository) QueryBuilder(doc map[string]any) (any, *options.FindOptions) {
	var filter bson.D
	opts := options.FindOptions{}
	filter = append(filter, bson.E{Key: "deleteDate", Value: nil})

	order := -1
	sort := "_id"
	var skip int
	var limit int
	project := bson.M{"_id": 0}
	for k, v := range doc {
		if k == "sort" {
			sort = v.(string)
		} else if k == "order" {
			if v.(string) == "asc" {
				order = 1
			}
		} else if k == "limit" {
			limit = v.(int)
		} else if k == "name" {
			filter = append(filter, bson.E{Key: k, Value: v})
		} else if k == "skip" {
			skip = v.(int)
		} else if k == "projection" {
			project = v.(bson.M)
		}

	}

	opts.SetProjection(project)
	opts.SetSort(bson.M{sort: order})
	fmt.Printf("sort: %s, order: %d\n", sort, order)
	opts.SetSkip(int64(skip))
	opts.SetLimit(int64(limit))

	return filter, &opts
}

func (p *productRepository) FindAll(doc map[string]any) ([]model.Products, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := []model.Products{}
	filter, opts := p.QueryBuilder(doc)

	cur, err := p.db.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var doc model.Products
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}

		doc.Href = utils.HostName(doc.Type, doc.ID)
		result = append(result, doc)
	}
	return result, nil
}

func (p *productRepository) CountProduct() int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := p.db.CountDocuments(ctx, bson.M{"deleteDate": nil})
	if err != nil {
		return 0
	}
	return count
}

func (p *productRepository) FindOne(id string) (*model.Products, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"deleteDate": nil,
		"id":         id,
	}

	result := model.Products{}

	if err := p.db.FindOne(ctx, filter).Decode(&result); err != nil {
		return nil, err
	}

	result.Href = utils.HostName(result.Type, result.ID)

	return &result, nil
}

func (p *productRepository) FindProductLanguage(id, languageCode string) (*model.ProductLanguage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "id", Value: id},
			{Key: "deleteDate", Value: primitive.Null{}},
		},
		}},
		bson.D{{Key: "$unwind", Value: "$supportingLanguage"}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "supportingLanguage.languageCode", Value: languageCode}}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "supportingLanguage", Value: bson.D{{Key: "$push", Value: "$supportingLanguage"}}},
		},
		}},
		bson.D{{Key: "$unwind", Value: "$supportingLanguage"}},
		bson.D{{Key: "$replaceWith", Value: "$supportingLanguage"}},
	}

	cur, err := p.db.Aggregate(ctx, filter)
	if err != nil {
		return nil, err
	}

	result := model.ProductLanguage{}
	for cur.Next(ctx) {
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (p *productRepository) FindProductLanguages(languageCode string) ([]model.ProductLanguage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "deleteDate", Value: primitive.Null{}}},
		}},
		bson.D{{Key: "$unwind", Value: "$supportingLanguage"}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "supportingLanguage.languageCode", Value: languageCode}}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "supportingLanguage", Value: bson.D{{Key: "$push", Value: "$supportingLanguage"}}},
		},
		}},
		bson.D{{Key: "$unwind", Value: "$supportingLanguage"}},
		bson.D{{Key: "$replaceWith", Value: "$supportingLanguage"}},
		bson.D{{Key: "$limit", Value: 20}},
	}

	productLanguages := []model.ProductLanguage{}
	cur, err := p.db.Aggregate(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		result := model.ProductLanguage{}
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}
		productLanguages = append(productLanguages, result)
	}

	return productLanguages, nil
}
