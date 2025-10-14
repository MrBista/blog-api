package dto

type UserResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name" validate:"required, min=1,max=100"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required, email"`
	Bio      string `json:"bio"`
	Role     int    `json:"role"`
}

type UserRequest struct {
	Name     string `json:"name" validate:"required, min=1,max=100"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password"`
}
