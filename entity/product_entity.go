package entity

import (
	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=2,max=100"`
	Price       float64   `gorm:"not null" json:"price" validate:"required,gt=0"`
	Description string    `gorm:"type:varchar(255)" json:"description" validate:"omitempty,max=255"`
	Quantity    int       `gorm:"not null" json:"quantity" validate:"required,min=0"`

	Timestamp
}
