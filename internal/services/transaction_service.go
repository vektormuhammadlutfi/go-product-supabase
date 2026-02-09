package services

import (
	"errors"
	"fmt"
	"gocats/internal/models"
	"gocats/internal/repository"

	"gorm.io/gorm"
)

type TransactionService interface {
	Checkout(request models.CheckoutRequest) (*models.Transaction, error)
	GetTransactionByID(id uint) (*models.Transaction, error)
	GetTodaySalesSummary() (*models.SalesSummary, error)
	GetSalesSummaryByDateRange(startDate, endDate string) (*models.SalesSummary, error)
}

type transactionService struct {
	db          *gorm.DB
	transRepo   repository.TransactionRepository
	productRepo repository.ProductRepository
}

func NewTransactionService(
	db *gorm.DB,
	transRepo repository.TransactionRepository,
	productRepo repository.ProductRepository) TransactionService {
	return &transactionService{
		db:          db,
		transRepo:   transRepo,
		productRepo: productRepo,
	}
}

func (s *transactionService) Checkout(request models.CheckoutRequest) (*models.Transaction, error) {
	if len(request.Items) == 0 {
		return nil, errors.New("checkout items cannot be empty")
	}

	var transaction *models.Transaction
	var transactionDetails []models.TransactionDetail

	// Start database transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var totalAmount float64

		// Validate all items and calculate total
		for _, item := range request.Items {
			if item.Quantity <= 0 {
				return fmt.Errorf("invalid quantity for product ID %d", item.ProductID)
			}

			// Get product and check stock using tx context
			var product models.Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				return fmt.Errorf("product ID %d not found", item.ProductID)
			}

			if product.Stock < item.Quantity {
				return fmt.Errorf("insufficient stock for product %s. Available: %d, Requested: %d",
					product.Name, product.Stock, item.Quantity)
			}

			subtotal := product.Price * float64(item.Quantity)
			totalAmount += subtotal

			// Store transaction detail for later creation
			transactionDetails = append(transactionDetails, models.TransactionDetail{
				ProductID: uint(item.ProductID),
				Quantity:  item.Quantity,
				Subtotal:  subtotal,
			})
		}

		// Create transaction
		transaction = &models.Transaction{
			TotalAmount: totalAmount,
		}

		if err := s.transRepo.CreateTransaction(tx, transaction); err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		// Create transaction details and update stock
		for i := range transactionDetails {
			transactionDetails[i].TransactionID = transaction.ID

			if err := s.transRepo.CreateTransactionDetail(tx, &transactionDetails[i]); err != nil {
				return fmt.Errorf("failed to create transaction detail: %w", err)
			}

			// Update product stock
			if err := s.transRepo.UpdateProductStock(tx, transactionDetails[i].ProductID, transactionDetails[i].Quantity); err != nil {
				return fmt.Errorf("failed to update product stock: %w", err)
			}
		}

		transaction.TransactionDetails = transactionDetails
		return nil
	})

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *transactionService) GetTransactionByID(id uint) (*models.Transaction, error) {
	transaction, err := s.transRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("transaction not found")
	}
	return transaction, nil
}

func (s *transactionService) GetTodaySalesSummary() (*models.SalesSummary, error) {
	summary, err := s.transRepo.GetTodaySummary()
	if err != nil {
		return nil, err
	}
	return summary, nil
}

func (s *transactionService) GetSalesSummaryByDateRange(startDate, endDate string) (*models.SalesSummary, error) {
	if startDate == "" || endDate == "" {
		return nil, errors.New("start_date and end_date are required")
	}

	summary, err := s.transRepo.GetSummaryByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}
	return summary, nil
}
