package models

/*
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}
*/

type Product struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	Name       string   `gorm:"size:200;not null" json:"name"`
	Price      float64  `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock      int      `gorm:"default:0" json:"stock"`
	CategoryID uint     `gorm:"not null;index" json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category"`
}

func (Product) TableName() string {
	return "products"
}

type ProductResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	Stock      int       `json:"stock"`
	CategoryID uint      `json:"category_id"`
	Category   *Category `json:"category,omitempty"`
}
