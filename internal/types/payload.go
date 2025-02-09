package types

type CreateUserPayload struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Email    string   `json:"email"`
	Role     UserRole `json:"role"`
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type CreateCompanyPayload struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type UpdateCompanyPayload struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type CreateCarPayload struct {
	Make           string  `json:"make"`
	Model          string  `json:"model"`
	Year           int     `json:"year"`
	Color          string  `json:"color"`
	RegistrationNo string  `json:"registration_no"`
	PricePerDay    float64 `json:"price_per_day"`
	CompanyID      int     `json:"company_id"`
}

type UpdateCarPayload struct {
	Make           string  `json:"make"`
	Model          string  `json:"model"`
	Year           int     `json:"year"`
	Color          string  `json:"color"`
	RegistrationNo string  `json:"registration_no"`
	PricePerDay    float64 `json:"price_per_day"`
}
