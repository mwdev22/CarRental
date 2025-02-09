package types

type UserRole int

const (
	UserTypeAdmin UserRole = iota
	UserTypeCompanyOwner
	UserTypeUser
)
