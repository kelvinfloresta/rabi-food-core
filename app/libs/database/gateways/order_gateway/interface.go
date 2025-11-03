package order_gateway

import (
	"rabi-food-core/domain"
	"rabi-food-core/libs/database"
	"time"
)

type OrderGateway interface {
	Create(input CreateInput) (string, error)
	GetByID(filter GetByIDFilter) (*GetByIDOutput, error)
	Patch(filter PatchFilter, values PatchValues) (bool, error)
	Paginate(filter PaginateFilter, paginate database.PaginateInput) (PaginateOutput, error)
	Delete(filter DeleteFilter) (bool, error)
}

type OrderItem struct {
	ProductID   string `json:"productId"`
	ProductName string `json:"productName"`
	Quantity    uint   `json:"quantity"`
	UnitPrice   uint   `json:"unitPrice"`
	Total       uint   `json:"total"`
}

type CreateInput struct {
	UserID     string
	TenantID   string
	Code       string
	Status     domain.OrderStatus
	Notes      string
	TotalPrice uint
	Items      []OrderItem
}

type GetByIDFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
}

type GetByIDOutput struct {
	ID         string             `json:"id"`
	TenantID   string             `json:"tenantId"`
	Code       string             `json:"code"`
	Status     domain.OrderStatus `json:"status"`
	Notes      string             `json:"notes"`
	TotalPrice uint               `json:"totalPrice"`
	Items      []OrderItem        `json:"items"`
	CreatedAt  time.Time          `json:"createdAt"`
}

type PatchFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
}

type PatchValues struct {
	Status domain.OrderStatus `json:"status"`
	Notes  string             `json:"notes"`
}

type PaginateFilter struct {
	UserID        *string            `json:"userId"`
	TenantID      *string            `json:"tenantId"`
	Status        domain.OrderStatus `json:"status"`
	CreatedAtFrom *time.Time         `json:"createdAtFrom"`
	CreatedAtTo   *time.Time         `json:"createdAtTo"`
}

type PaginateData struct {
	ID         string             `json:"id"`
	TenantID   string             `json:"tenantId"`
	Code       string             `json:"code"`
	Status     domain.OrderStatus `json:"status"`
	Notes      string             `json:"notes"`
	TotalPrice uint               `json:"totalPrice"`
	CreatedAt  time.Time          `json:"createdAt"`
}

type PaginateOutput struct {
	Data     []PaginateData `json:"data"`
	MaxPages int            `json:"maxPages"`
}

type DeleteFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
}
