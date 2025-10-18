package dto

import "time"

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

type UserFollowerDTO struct {
	ID              uint64    `json:"id"`
	Name            string    `json:"name"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	ProfileImageURI *string   `json:"profile_image_uri"`
	Bio             *string   `json:"bio"`
	FollowedAt      time.Time `json:"followed_at"`
}

type UserFollowingDTO struct {
	ID              uint64    `json:"id"`
	Name            string    `json:"name"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	ProfileImageURI *string   `json:"profile_image_uri"`
	Bio             *string   `json:"bio"`
	FollowedAt      time.Time `json:"followed_at"`
}
