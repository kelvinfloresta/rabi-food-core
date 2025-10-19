package models

type Tenant struct {
	ID   string `gorm:"type:uuid"`
	Name string `gorm:"not null"`
}

func (m Tenant) TableName() string {
	return "tenants"
}
