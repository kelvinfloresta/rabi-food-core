package category_gateway

import (
	"rabi-food-core/libs/database"
)

type CategoryGateway interface {
	Create(input CreateInput) (string, error)
	GetByID(id string, tenantID string) (*GetByIDOutput, error)
	Patch(filter PatchFilter, values PatchValues) (bool, error)
	Paginate(filter PaginateFilter, paginate database.PaginateInput) (PaginateOutput, error)
	Delete(id string) (bool, error)
}

type CreateInput struct {
	TenantID    string `json:"tenantId"`
	Name        string `validate:"required" json:"name"`
	Description string `json:"description"`
}

type GetByIDOutput struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TenantID    string `json:"tenantId"`
}

type PatchFilter struct {
	ID string
}

type PatchValues struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PaginateFilter struct {
	Name *string `json:"name"`
}

type PaginateData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PaginateOutput struct {
	Data     []PaginateData `json:"data"`
	MaxPages int            `json:"maxPages"`
}
