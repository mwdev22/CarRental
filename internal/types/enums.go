package types

type UserRole int

const (
	UserTypeTest UserRole = iota
	UserTypeAdmin
	UserTypeCompanyOwner
	UserTypeUser
)
