package dto

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
	Type       string `json:"type"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	TypeToken    string `json:"typeToken"`
	RefreshToken string `json:"refreshToken"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Bio      string `json:"bio"`
	Role     int    `json:"role"`
}
