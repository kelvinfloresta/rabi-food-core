package domain

type Role string

const (
	UserRole       Role = "user"
	BackofficeRole Role = "backoffice"
)

func (r Role) IsBackoffice() bool {
	return r == BackofficeRole
}

func (r Role) IsUser() bool {
	return r == UserRole
}
