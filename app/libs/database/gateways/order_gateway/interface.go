package order_gateway

import (
	"rabi-food-core/domain/order"
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
	UserID            string
	TenantID          string
	Code              string
	PaymentStatus     order.PaymentStatus
	FulfillmentStatus order.FulfillmentStatus
	DeliveryStatus    order.DeliveryStatus
	Notes             string
	TotalPrice        uint
	Items             []OrderItem
}

type GetByIDFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
}

type GetByIDOutput struct {
	ID                string                  `json:"id"`
	TenantID          string                  `json:"tenantId"`
	Code              string                  `json:"code"`
	PaymentStatus     order.PaymentStatus     `json:"paymentStatus"`
	FulfillmentStatus order.FulfillmentStatus `json:"fulfillmentStatus"`
	DeliveryStatus    order.DeliveryStatus    `json:"deliveryStatus"`
	Notes             string                  `json:"notes"`
	TotalPrice        uint                    `json:"totalPrice"`
	Items             []OrderItem             `json:"items"`
	CreatedAt         time.Time               `json:"createdAt"`
}

type PatchFilter struct {
	ID                  string
	TenantID            string
	DeliveryStatusIn    []order.DeliveryStatus
	FulfillmentStatusIn []order.FulfillmentStatus
	PaymentStatusIn     []order.PaymentStatus
}

type PatchValues struct {
	PaymentStatus     order.PaymentStatus     `json:"paymentStatus"`
	FulfillmentStatus order.FulfillmentStatus `json:"fulfillmentStatus"`
	DeliveryStatus    order.DeliveryStatus    `json:"deliveryStatus"`
}

type PaginateFilter struct {
	UserID            *string                 `json:"userId"`
	TenantID          *string                 `json:"tenantId"`
	PaymentStatus     order.PaymentStatus     `json:"paymentStatus"`
	FulfillmentStatus order.FulfillmentStatus `json:"fulfillmentStatus"`
	DeliveryStatus    order.DeliveryStatus    `json:"deliveryStatus"`
	CreatedAtFrom     *time.Time              `json:"createdAtFrom"`
	CreatedAtTo       *time.Time              `json:"createdAtTo"`
}

type PaginateData struct {
	ID                string                  `json:"id"`
	TenantID          string                  `json:"tenantId"`
	Code              string                  `json:"code"`
	PaymentStatus     order.PaymentStatus     `json:"paymentStatus"`
	FulfillmentStatus order.FulfillmentStatus `json:"fulfillmentStatus"`
	DeliveryStatus    order.DeliveryStatus    `json:"deliveryStatus"`
	Notes             string                  `json:"notes"`
	TotalPrice        uint                    `json:"totalPrice"`
	CreatedAt         time.Time               `json:"createdAt"`
}

type PaginateOutput struct {
	Data     []PaginateData `json:"data"`
	MaxPages int            `json:"maxPages"`
}

type DeleteFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
}
