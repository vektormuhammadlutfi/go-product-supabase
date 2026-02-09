package models

type TransactionDetail struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	TransactionID uint    `gorm:"not null;index" json:"transaction_id"`
	ProductID     uint    `gorm:"not null;index" json:"product_id"`
	Quantity      int     `gorm:"not null" json:"quantity"`
	Subtotal      float64 `gorm:"type:decimal(10,2);not null" json:"subtotal"`

	Transaction Transaction `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
	Product     Product     `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (TransactionDetail) TableName() string {
	return "transaction_details"
}
