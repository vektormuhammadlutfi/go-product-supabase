package repository

import (
	"gocats/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(tx *gorm.DB, transaction *models.Transaction) error
	CreateTransactionDetail(tx *gorm.DB, detail *models.TransactionDetail) error
	UpdateProductStock(tx *gorm.DB, productID uint, quantity int) error
	FindByID(id uint) (*models.Transaction, error)
	GetTodaySummary() (*models.SalesSummary, error)
	GetSummaryByDateRange(startDate, endDate string) (*models.SalesSummary, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(tx *gorm.DB, transaction *models.Transaction) error {
	return tx.Create(transaction).Error
}

func (r *transactionRepository) CreateTransactionDetail(tx *gorm.DB, detail *models.TransactionDetail) error {
	return tx.Create(detail).Error
}

func (r *transactionRepository) UpdateProductStock(tx *gorm.DB, productID uint, quantity int) error {
	return tx.Model(&models.Product{}).Where("id = ?", productID).UpdateColumn("stock", gorm.Expr("stock - ?", quantity)).Error
}

func (r *transactionRepository) FindByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("TransactionDetails.Product").First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) GetTodaySummary() (*models.SalesSummary, error) {
	var summary models.SalesSummary

	// Get total revenue and total transactions for today
	type Result struct {
		TotalRevenue      float64
		TotalTransactions int
	}
	var result Result

	err := r.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(total_amount), 0) as total_revenue, COUNT(id) as total_transactions").
		Where("DATE(created_at) = CURRENT_DATE").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	summary.TotalRevenue = result.TotalRevenue
	summary.TotalTransactions = result.TotalTransactions

	// Get best selling product for today
	type BestProduct struct {
		ProductID uint
		Name      string
		QtySold   int
	}

	var bestProduct BestProduct
	err = r.db.Model(&models.TransactionDetail{}).
		Select("transaction_details.product_id, products.name, SUM(transaction_details.quantity) as qty_sold").
		Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").
		Joins("JOIN products ON products.id = transaction_details.product_id").
		Where("DATE(transactions.created_at) = CURRENT_DATE").
		Group("transaction_details.product_id, products.name").
		Order("qty_sold DESC").
		Limit(1).
		Scan(&bestProduct).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Set best selling product if found
	if bestProduct.ProductID != 0 {
		summary.BestSellingProduct = &models.BestSellingProduct{
			Name:    bestProduct.Name,
			QtySold: bestProduct.QtySold,
		}
	}

	return &summary, nil
}

func (r *transactionRepository) GetSummaryByDateRange(startDate, endDate string) (*models.SalesSummary, error) {
	var summary models.SalesSummary

	// Get total revenue and total transactions for date range
	type Result struct {
		TotalRevenue      float64
		TotalTransactions int
	}
	var result Result

	err := r.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(total_amount), 0) as total_revenue, COUNT(id) as total_transactions").
		Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", startDate, endDate).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	summary.TotalRevenue = result.TotalRevenue
	summary.TotalTransactions = result.TotalTransactions

	// Get best selling product for date range
	type BestProduct struct {
		ProductID uint
		Name      string
		QtySold   int
	}

	var bestProduct BestProduct
	err = r.db.Model(&models.TransactionDetail{}).
		Select("transaction_details.product_id, products.name, SUM(transaction_details.quantity) as qty_sold").
		Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").
		Joins("JOIN products ON products.id = transaction_details.product_id").
		Where("DATE(transactions.created_at) >= ? AND DATE(transactions.created_at) <= ?", startDate, endDate).
		Group("transaction_details.product_id, products.name").
		Order("qty_sold DESC").
		Limit(1).
		Scan(&bestProduct).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Set best selling product if found
	if bestProduct.ProductID != 0 {
		summary.BestSellingProduct = &models.BestSellingProduct{
			Name:    bestProduct.Name,
			QtySold: bestProduct.QtySold,
		}
	}

	return &summary, nil
}
