package types

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
