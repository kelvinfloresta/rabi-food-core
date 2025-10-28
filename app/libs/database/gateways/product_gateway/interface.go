package product_gateway

import (
	"rabi-food-core/libs/database"
)

type ProductGateway interface {
	Create(input CreateInput) (string, error)
	GetByID(id string) (*GetByIDOutput, error)
	Patch(filter PatchFilter, values PatchValues) (bool, error)
	Paginate(filter PaginateFilter, paginate database.PaginateInput) (PaginateOutput, error)
	Delete(id string) (bool, error)
}

type CreateInput struct {
	TenantID    string `json:"tenantId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
	CategoryID  string `json:"categoryId"`
	Unit        string `json:"unit"`
	Price       int    `json:"price"`
	IsActive    bool   `json:"isActive"`
}

type GetByIDOutput struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Photo        string `json:"photo"`
	CategoryID   string `json:"categoryId"`
	CategoryName string `json:"categoryName"`
	Unit         string `json:"unit"`
	Price        int    `json:"price"`
	IsActive     bool   `json:"isActive"`
}

type PatchFilter struct {
	ID string
}

type PatchValues struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
	CategoryID  string `json:"categoryId"`
	Unit        string `json:"unit"`
	Price       int    `json:"price"`
	IsActive    bool   `json:"isActive"`
}

type PaginateFilter struct {
	Name       *string `json:"name"`
	CategoryID *string `json:"categoryId"`
	IsActive   *bool   `json:"isActive"`
}

type PaginateData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
	CategoryID  string `json:"categoryId"`
	Unit        string `json:"unit"`
	Price       int    `json:"price"`
	IsActive    bool   `json:"isActive"`
}

type PaginateOutput struct {
	Data     []PaginateData `json:"data"`
	MaxPages int            `json:"maxPages"`
}
