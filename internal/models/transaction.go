package models

import "time"

type Transaction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TotalAmount float64   `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	TransactionDetails []TransactionDetail `gorm:"foreignKey:TransactionID;constraint:OnDelete:CASCADE" json:"transaction_details,omitempty"`
}

func (Transaction) TableName() string {
	return "transactions"
}

// CheckoutItem represents a single item in the checkout request
type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// CheckoutRequest represents the checkout request payload
type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

// BestSellingProduct represents the best selling product info
type BestSellingProduct struct {
	Name    string `json:"name"`
	QtySold int    `json:"qty_sold"`
}

// SalesSummary represents the sales summary response
type SalesSummary struct {
	TotalRevenue       float64             `gorm:"-" json:"total_revenue"`
	TotalTransactions  int                 `gorm:"-" json:"total_transactions"`
	BestSellingProduct *BestSellingProduct `gorm:"-" json:"best_selling_product,omitempty"`
}
