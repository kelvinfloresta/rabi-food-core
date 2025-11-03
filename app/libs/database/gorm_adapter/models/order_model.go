package models

import (
	"rabi-food-core/domain"
	"time"

	"gorm.io/datatypes"
)

// Order represents the Order model in the database.
type Order struct {
	ID       string `gorm:"type:uuid"`
	TenantID string `gorm:"type:uuid;not null"`
	Tenant   Tenant
	UserID   string `gorm:"type:uuid;not null"`
	User     User

	Code       string             `gorm:"uniqueIndex;not null"`
	Status     domain.OrderStatus `gorm:"type:varchar(20);not null"`
	Notes      string             `gorm:"type:text"`
	TotalPrice uint               `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`

	Items datatypes.JSON `gorm:"not null"`
}

func (Order) TableName() string {
	return "orders"
}
