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

type GoogleCallbackRequest struct {
	Code string `json:"code" validate:"required"`
}

type GoogleAuthURLResponse struct {
	URL string `json:"url"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verifiedEmail"`
	Name          string `json:"name"`
	GivenName     string `json:"givenName"`
	FamilyName    string `json:"familyName"`
	Picture       string `json:"picture"`
}
