package types

type CreateUserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type CreateCarPayload struct {
	Make           string  `json:"make"`
	Model          string  `json:"model"`
	Year           int     `json:"year"`
	Color          string  `json:"color"`
	RegistrationNo string  `json:"registration_no"`
	PricePerDay    float64 `json:"price_per_day"`
}

type UpdateCarPayload struct {
	Make           string  `json:"make"`
	Model          string  `json:"model"`
	Year           int     `json:"year"`
	Color          string  `json:"color"`
	RegistrationNo string  `json:"registration_no"`
	PricePerDay    float64 `json:"price_per_day"`
}
