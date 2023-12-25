package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/aidarkhanov/nanoid/v2"
	"github.com/bxcodec/faker/v3"
	"github.com/sing3demons/product/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Seed struct {
	db *mongo.Client
}

const (
	alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func RandomNanoID(size int) string {
	id, _ := nanoid.GenerateString(alphabet, size)
	return id
}

func InitSeed(db *mongo.Client) {
	seed := Seed{db: db}

	categoryDB := db.Database("category").Collection("category")
	languageDB := db.Database("language").Collection("productLanguage")
	productDB := db.Database("products").Collection("products")

	optionsIndex := options.Index()
	optionsIndex.SetUnique(true)
	optionsIndex.SetSparse(true)

	categoryDB.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1, "name": 1},
			Options: optionsIndex,
		},
	})
	productDB.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: optionsIndex,
		},
	})
	languageDB.Indexes().CreateMany(context.Background(), []mongo.IndexModel{{
		Keys:    bson.M{"id": 1},
		Options: optionsIndex,
	}})

	categories := []Category{}

	opt := options.Find()
	opt.SetSort(bson.M{"_id": -1})
	cur, err := categoryDB.Find(context.Background(), bson.M{"deleteDate": nil}, opt)
	if err != nil {
		panic(err)
	}
	if err := cur.All(context.Background(), &categories); err != nil {
		panic(err)
	}

	fmt.Println("Category already created :: ", len(categories))

	if len(categories) == 0 {
		seed.createCategory(categoryDB)
		return
	} else {
		start := time.Now()
		for i := 0; i < 9312; i++ {
			var randomCategory []model.Category
			random := rand.Intn(len(categories))
			for i := 0; i < random; i++ {
				result := model.Category{
					ID:   categories[i].ID,
					Name: categories[i].Name,
				}
				randomCategory = append(randomCategory, result)
			}

			supportingLanguage := []model.ProductLanguage{}

			documents := []interface{}{
				model.ProductLanguage{
					Type:         "productLanguage",
					ID:           RandomNanoID(11),
					LanguageCode: "en",
					Name:         faker.Name(),
					Version:      "1.0",
					LastUpdate:   time.Now().Format(time.RFC3339),
				},
				model.ProductLanguage{
					Type:         "productLanguage",
					ID:           RandomNanoID(11),
					LanguageCode: "ru",
					Name:         faker.Name(),
					Version:      "1.0",
					LastUpdate:   time.Now().Format(time.RFC3339),
				},
				model.ProductLanguage{
					Type:         "productLanguage",
					ID:           RandomNanoID(11),
					LanguageCode: "th",
					Name:         faker.Name(),
					Version:      "1.0",
					LastUpdate:   time.Now().Format(time.RFC3339),
				},
				model.ProductLanguage{
					Type:         "productLanguage",
					ID:           RandomNanoID(11),
					LanguageCode: "zh",
					Name:         faker.Name(),
					Version:      "1.0",
					LastUpdate:   time.Now().Format(time.RFC3339),
				},
				model.ProductLanguage{
					Type:         "productLanguage",
					ID:           RandomNanoID(11),
					LanguageCode: "ja",
					Name:         faker.Name(),
					Version:      "1.0",
					LastUpdate:   time.Now().Format(time.RFC3339)},
				model.ProductLanguage{Type: "productLanguage",
					ID:           RandomNanoID(11),
					LanguageCode: "ko",
					Name:         faker.Name(),
					Version:      "1.0",
					LastUpdate:   time.Now().Format(time.RFC3339)},
			}

			languageDB.InsertMany(context.Background(), documents)

			for _, lang := range documents {
				result := lang.(model.ProductLanguage)
				supportingLanguage = append(supportingLanguage, model.ProductLanguage{
					ID:           result.ID,
					LanguageCode: result.LanguageCode,
					Name:         result.Name,
					Type:         result.Type,
				})
			}

			product := model.Products{
				Type:            "products",
				ID:              RandomNanoID(11),
				Name:            faker.Name(),
				Version:         "1.0",
				LastUpdate:      time.Now().Format(time.RFC3339),
				LifecycleStatus: "active",
				Category:        randomCategory,
				ValidFor: &model.ValidFor{
					StartDateTime: time.Now().UTC(),
					EndDateTime:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
				},
				SupportingLanguage: supportingLanguage,
			}

			seed.SeedProducts(productDB, product)

			for _, cat := range product.Category {
				update := bson.M{"$push": bson.M{"products": bson.M{"id": product.ID, "name": product.Name}}}
				categoryDB.UpdateOne(context.Background(), bson.M{"id": cat.ID}, update)
			}
			fmt.Printf("Product created :: %d\n", i)
		}
		fmt.Println("Time :: ", time.Since(start))
	}
}

func (s *Seed) SeedProducts(db *mongo.Collection, product any) {
	_, err := db.InsertOne(context.Background(), product)
	if err != nil {
		panic(err)
	}
	// fmt.Println(result.InsertedID)
}

func (s *Seed) createCategory(db *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	categories := []string{
		"Mobile",
		"Fixed",
		"TV",
		"Internet",
		"Home",
		"Business",
		"Entertainment",
		"Sport",
		"Music",
		"News",
		"Kids",
		"Education",
		"Health",
		"Travel",
		"Finance",
		"Shopping",
		"Security",
		"Utilities",
		"Transport",
		"Other",
	}

	for _, name := range categories {
		category := Category{
			Type:            "category",
			ID:              RandomNanoID(11),
			Name:            name,
			Version:         "1.0",
			LastUpdate:      time.Now().Format(time.RFC3339),
			LifecycleStatus: "active",
		}
		result, err := db.InsertOne(ctx, category)
		if err != nil {
			panic(err)
		}
		fmt.Println(result.InsertedID)
	}

	total, err := db.CountDocuments(ctx, bson.M{"deleteDate": nil})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Category created :: ", total)

}

type Category struct {
	Type            string    `json:"@type,omitempty" validate:"required" bson:"@type,omitempty"`
	ID              string    `json:"id" validate:"required" bson:"id"`
	Href            string    `json:"href,omitempty" bson:"href,omitempty"`
	Name            string    `json:"name,omitempty" bson:"name,omitempty"`
	Version         string    `json:"version,omitempty" bson:"version,omitempty"`
	LastUpdate      string    `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	ValidFor        *ValidFor `json:"validFor,omitempty" bson:"validFor,omitempty"`
	LifecycleStatus string    `json:"lifecycleStatus,omitempty" bson:"lifecycleStatus,omitempty"`
}

type ValidFor struct {
	StartDateTime string `json:"startDateTime,omitempty" bson:"startDateTime,omitempty"`
	EndDateTime   string `json:"endDateTime,omitempty" bson:"endDateTime,omitempty"`
}
