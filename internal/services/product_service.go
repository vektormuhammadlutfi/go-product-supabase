package services

import (
	"errors"
	"gocats/internal/models"
	"gocats/internal/repository"

	"gorm.io/gorm"
)

type ProductService interface {
	CreateProduct(name, description string, price float64, stock int, categoryID uint) (*models.Product, error)
	GetAllProducts(name string) ([]models.ProductResponse, error)
	GetProductByID(id uint) (*models.ProductResponse, error)
	GetProductsByCategoryID(categoryID uint) ([]models.ProductResponse, error)
	UpdateProduct(id uint, name, description string, price float64, stock int, categoryID uint) (*models.Product, error)
	DeleteProduct(id uint) error
}

type productService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

func NewProductService(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository) ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *productService) CreateProduct(name, description string, price float64, stock int, categoryID uint) (*models.Product, error) {
	// Implementation goes here
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}

	if price < 0 {
		return nil, errors.New("product price cannot be negative")
	}

	if stock < 0 {
		return nil, errors.New("product stock cannot be negative")
	}

	_, err := s.categoryRepo.FindByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	product := &models.Product{
		Name:       name,
		Price:      price,
		Stock:      stock,
		CategoryID: categoryID,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) GetAllProducts(name string) ([]models.ProductResponse, error) {
	var products []models.Product
	var err error

	if name != "" {
		products, err = s.productRepo.FindByName(name)
	} else {
		products, err = s.productRepo.FindAll()
	}

	if err != nil {
		return nil, err
	}

	responses := make([]models.ProductResponse, len(products))
	for i, product := range products {
		var cat *models.Category
		if product.Category.ID != 0 || product.Category.Name != "" {
			cat = &models.Category{
				ID:          product.Category.ID,
				Name:        product.Category.Name,
				Description: product.Category.Description,
			}
		}
		responses[i] = models.ProductResponse{
			ID:         product.ID,
			Name:       product.Name,
			Price:      product.Price,
			Stock:      product.Stock,
			CategoryID: product.CategoryID,
			Category:   cat,
		}
	}

	return responses, nil
}

func (s *productService) GetProductByID(id uint) (*models.ProductResponse, error) {
	if id == 0 {
		return nil, errors.New("product ID cannot be zero")
	}

	product, err := s.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	var cat *models.Category
	if product.Category.ID != 0 || product.Category.Name != "" {
		cat = &models.Category{
			ID:          product.Category.ID,
			Name:        product.Category.Name,
			Description: product.Category.Description,
		}
	}

	response := &models.ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		Price:      product.Price,
		Stock:      product.Stock,
		CategoryID: product.CategoryID,
		Category:   cat,
	}

	return response, nil
}

func (s *productService) GetProductsByCategoryID(categoryID uint) ([]models.ProductResponse, error) {
	if categoryID == 0 {
		return nil, errors.New("category ID cannot be zero")
	}

	_, err := s.categoryRepo.FindByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	products, err := s.productRepo.FindByCategoryID(categoryID)
	if err != nil {
		return nil, err
	}

	responses := make([]models.ProductResponse, len(products))
	for i, product := range products {
		var cat *models.Category
		if product.Category.ID != 0 || product.Category.Name != "" {
			cat = &models.Category{
				ID:          product.Category.ID,
				Name:        product.Category.Name,
				Description: product.Category.Description,
			}
		}
		responses[i] = models.ProductResponse{
			ID:         product.ID,
			Name:       product.Name,
			Price:      product.Price,
			Stock:      product.Stock,
			CategoryID: product.CategoryID,
			Category:   cat,
		}
	}

	return responses, nil
}

func (s *productService) UpdateProduct(id uint, name, description string, price float64, stock int, categoryID uint) (*models.Product, error) {
	if id == 0 {
		return nil, errors.New("product ID cannot be zero")
	}

	product, err := s.productRepo.FindByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
	}

	if name != "" {
		product.Name = name
	}

	if stock >= 0 {
		product.Stock = stock
	}

	if price >= 0 {
		product.Price = price
	}

	if stock >= 0 {
		product.Stock = stock
	}

	if categoryID > 0 {
		_, err := s.categoryRepo.FindByID(categoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("category not found")
			}
			return nil, err
		}
		product.CategoryID = categoryID
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	return product, nil

}

func (s *productService) DeleteProduct(id uint) error {
	if id == 0 {
		return errors.New("product ID cannot be zero")
	}

	_, err := s.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	return s.productRepo.Delete(id)
}
