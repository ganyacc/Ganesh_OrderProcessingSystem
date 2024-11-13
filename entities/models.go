package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderStatus represents the status of an order.
type OrderStatus string

// Define constants for the order status values.
const (
	Unfulfilled OrderStatus = "unfulfilled"
	Fulfilled   OrderStatus = "fulfilled"
)

type BaseModel struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return
}

type Customer struct {
	BaseModel
	Name    string  `json:"name" validate:"required"`
	Email   string  `json:"email" validate:"required"`
	Country string  `json:"country"`
	Order   []Order `gorm:"foreignKey:CustomerID" json:"orders"`
}

type Product struct {
	BaseModel
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
}

type Order struct {
	BaseModel
	CustomerID uuid.UUID   `json:"customer_id"`
	Customer   Customer    `gorm:"foreignKey:CustomerID" json:"-"`
	Products   []Product   `gorm:"many2many:order_products;" json:"products"`
	TotalPrice float64     `json:"total_price"`
	Status     OrderStatus `json:"status" validate:"required,oneof=fulfilled unfulfilled"`
}

type OrderRequest struct {
	CustomerID string   `json:"customer_id" validate:"required"`
	ProductIDs []string `json:"product_ids" validate:"required"`
}
