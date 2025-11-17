package domain

type Role string

const (
	UserRole       Role = "user"
	StaffRole      Role = "staff"
	BackofficeRole Role = "backoffice"
)

func (r Role) IsBackoffice() bool {
	return r == BackofficeRole
}

func (r Role) IsUser() bool {
	return r == UserRole
}

func (r Role) IsStaff() bool {
	return r == StaffRole
}
