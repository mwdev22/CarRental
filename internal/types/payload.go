package types

// Validator instance

type CreateUserPayload struct {
	Username string   `json:"username" validate:"required,min=3,max=30"`
	Password string   `json:"password" validate:"required,min=3"`
	Email    string   `json:"email" validate:"required"`
	Role     UserRole `json:"role" validate:"required"`
}

type LoginPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserPayload struct {
	Username string `json:"username" validate:"omitempty,min=3,max=30"`
	Email    string `json:"email" validate:"omitempty"`
}

type CreateCompanyPayload struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Address string `json:"address" validate:"required"`
}

type UpdateCompanyPayload struct {
	Name    string `json:"name" validate:"omitempty"`
	Email   string `json:"email" validate:"omitempty"`
	Phone   string `json:"phone" validate:"omitempty"`
	Address string `json:"address" validate:"omitempty"`
}

type CreateCarPayload struct {
	Make           string  `json:"make" validate:"required"`
	Model          string  `json:"model" validate:"required"`
	Year           int     `json:"year" validate:"required,gte=1886,lte=2025"`
	Color          string  `json:"color" validate:"required"`
	RegistrationNo string  `json:"registration_no" validate:"required"`
	PricePerDay    float64 `json:"price_per_day" validate:"required,gt=0"`
	CompanyID      int     `json:"company_id" validate:"required"`
}

type UpdateCarPayload struct {
	Make           string  `json:"make" validate:"omitempty"`
	Model          string  `json:"model" validate:"omitempty"`
	Year           int     `json:"year" validate:"omitempty,gte=1886,lte=2025"`
	Color          string  `json:"color" validate:"omitempty"`
	RegistrationNo string  `json:"registration_no" validate:"omitempty"`
	PricePerDay    float64 `json:"price_per_day" validate:"omitempty,gt=0"`
}

type CreateBookingPayload struct {
	CarID     int    `json:"car_id" validate:"required"`
	StartDate string `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `json:"end_date" validate:"required,datetime=2006-01-02"`
}

type UpdateBookingPayload struct {
	StartDate string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate   string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
}
