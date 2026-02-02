package models

/*
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

*/

type Category struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}
