package services

import (
	"errors"
	"gocats/internal/models"
	"gocats/internal/repository"

	"gorm.io/gorm"
)

type CategoryService interface {
	CreateCategory(name, description string) (*models.Category, error)
	GetAllCategories() ([]models.Category, error)
	GetCategoryByID(id uint) (*models.Category, error)
	UpdateCategory(id uint, name, description string) (*models.Category, error)
	DeleteCategory(id uint) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(name, description string) (*models.Category, error) {
	if name == "" {
		return nil, errors.New("category name is required")
	}

	category := &models.Category{
		Name:        name,
		Description: description,
	}

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetAllCategories() ([]models.Category, error) {
	return s.repo.FindAll()
}

func (s *categoryService) GetCategoryByID(id uint) (*models.Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return category, nil
}

func (s *categoryService) UpdateCategory(id uint, name, description string) (*models.Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	if name != "" {
		category.Name = name
	}
	if description != "" {
		category.Description = description
	}

	if err := s.repo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}

	return s.repo.Delete(id)
}
