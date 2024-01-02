package service

import (
	"os"
	"sync"

	"github.com/sing3demons/product/model"
	"github.com/sing3demons/product/repository"
	"github.com/sing3demons/product/utils"
)

type ProductService interface {
	FindAll(doc map[string]any) ([]model.Products, int64, error)
	FindOne(id string) (*model.Products, error)
	FindProductLanguage(id, languageCode string) (*model.ProductLanguage, error)
	FindProductLanguages(languageCode string) ([]model.ProductLanguage, error)
}

type productService struct {
	repository      repository.ProductRepository
	categoryService CategoryService
}

func NewProductService(repository repository.ProductRepository, categoryService CategoryService) ProductService {
	return &productService{repository, categoryService}
}

func (svc *productService) FindAll(doc map[string]any) ([]model.Products, int64, error) {
	products, err := svc.repository.FindAll(doc)
	if err != nil {
		return nil, 0, err
	}

	if doc["language"] == true {
		for i := 0; i < len(products); i++ {
			if len(products[i].Category) != 0 {
				// products[i].Category = svc.GetCategory(products[i].Category)
				products[i].Category = svc.GetCategoryGrpc(products[i].Category)
			}

			if len(products[i].SupportingLanguage) != 0 {
				products[i].SupportingLanguage = svc.GetProductLanguage(products[i].SupportingLanguage)
			}
		}
	} else {
		for index := 0; index < len(products); index++ {
			// products[index].Category = svc.GetCategory(products[index].Category)
			products[index].Category = svc.GetCategoryGrpc(products[index].Category)
			products[index].SupportingLanguage = svc.GetProductLanguage(products[index].SupportingLanguage)
		}
	}

	total := svc.repository.CountProduct()

	return products, total, nil
}

func (svc *productService) FindOne(id string) (*model.Products, error) {
	product, err := svc.repository.FindOne(id)
	if err != nil {
		return nil, err
	}

	// if len(product.Category) != 0 {
	// 	product.Category = svc.GetCategory(product.Category)
	// }
	// if len(product.Category) != 0 {
	// 	result := []model.Category{}
	// 	for _, category := range product.Category {
	// 		r, err := svc.categoryService.GetCategory(category.ID)
	// 		if err != nil {
	// 			r = &category
	// 		}
	// 		result = append(result, *r)
	// 	}

	// 	product.Category = result
	// }

	if len(product.Category) != 0 {
		product.Category = svc.GetCategoryGrpc(product.Category)
	}

	if len(product.SupportingLanguage) != 0 {
		product.SupportingLanguage = svc.GetProductLanguage(product.SupportingLanguage)
	}

	return product, nil
}

func (svc *productService) GetCategoryGrpc(categories []model.Category) []model.Category {
	var wg sync.WaitGroup
	var newCategories []model.Category
	for _, category := range categories {
		wg.Add(1)
		go func(category model.Category) {
			defer wg.Done()
			c, err := svc.categoryService.GetCategory(category.ID)
			if err != nil {
				category.Type = "category"
				newCategories = append(newCategories, category)
				return
			}
			newCategories = append(newCategories, *c)
		}(category)
	}
	wg.Wait()
	return newCategories
}

func (svc *productService) FindProductLanguage(id, languageCode string) (*model.ProductLanguage, error) {
	productLanguage, err := svc.repository.FindProductLanguage(id, languageCode)
	if err != nil {
		return nil, err
	}

	url := utils.Href(os.Getenv("PRODUCT_LANGUAGE_SERVICE_URL"), productLanguage.Type, productLanguage.ID)
	result, err := utils.RequestHttpGet[model.ProductLanguage](url)
	if err != nil {
		result = productLanguage
	}

	return result, nil
}

func (svc *productService) GetCategory(categories []model.Category) []model.Category {
	var wg sync.WaitGroup
	host := os.Getenv("CATEGORY_SERVICE_URL")
	var newCategories []model.Category
	for _, category := range categories {
		wg.Add(1)
		go func(category model.Category) {
			defer wg.Done()
			url := host + "/" + category.Type + "/" + category.ID
			c, err := utils.RequestHttpGet[model.Category](url)
			if err != nil {
				category.Type = "category"
				newCategories = append(newCategories, category)
				return
			}
			newCategories = append(newCategories, *c)
		}(category)
	}
	wg.Wait()
	return newCategories
}

func (svc *productService) GetProductLanguage(supportingLanguage []model.ProductLanguage) []model.ProductLanguage {
	var wg sync.WaitGroup
	host := os.Getenv("PRODUCT_LANGUAGE_SERVICE_URL")
	var productLanguages []model.ProductLanguage
	for _, lang := range supportingLanguage {
		wg.Add(1)
		go func(lang model.ProductLanguage) {
			defer wg.Done()
			url := host + "/" + lang.Type + "/" + lang.ID
			productLanguage, err := utils.RequestHttpGet[model.ProductLanguage](url)
			if err != nil {
				productLanguages = append(productLanguages, lang)
				return
			}
			productLanguages = append(productLanguages, *productLanguage)
		}(lang)
	}
	wg.Wait()
	return productLanguages
}

func (svc *productService) FindProductLanguages(languageCode string) ([]model.ProductLanguage, error) {
	productLanguage, err := svc.repository.FindProductLanguages(languageCode)
	if err != nil {
		return nil, err
	}
	productLanguages := []model.ProductLanguage{}
	host := os.Getenv("PRODUCT_LANGUAGE_SERVICE_URL")
	for _, lang := range productLanguage {
		url := utils.Href(host, lang.Type, lang.ID)
		result, err := utils.RequestHttpGet[model.ProductLanguage](url)
		if err != nil {
			result = &lang
		}
		productLanguages = append(productLanguages, *result)
	}

	return productLanguages, nil
}
