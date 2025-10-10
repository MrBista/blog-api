package dto

import "time"

type CreatePostRequest struct {
	Id         int    `json:"id,omitempty"`
	Title      string `json:"title" validate:"required"`
	Content    string `json:"content" validate:"required"`
	Slug       string `json:"slug" validate:"required"`
	CategoryId int    `jsong:"categoryId" validate:"required"`
}

type UpdatePostRequest struct {
	Id      int     `json:"id,omitempty"`
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Status  int     `json:"status"`
}

type PostResponse struct {
	ID           uint64    `json:"id"`
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	Content      string    `json:"content"`
	MainImageURI string    `json:"mainImageURI"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
