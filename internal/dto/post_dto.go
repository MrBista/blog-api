package dto

import "time"

type CreatePostRequest struct {
	Id         int    `json:"id,omitempty"`
	Title      string `json:"title" validate:"required"`
	Content    string `json:"content" validate:"required"`
	CategoryId int    `json:"categoryId" validate:"required"`
}

type UpdatePostRequest struct {
	Slug    string  `json:"slug,omitempty"`
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Status  int     `json:"status" validate:"required"`
}

type AuthorResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Id    int64  `json:"id"`
}

type PostResponse struct {
	ID             uint64            `json:"id"`
	Title          string            `json:"title"`
	Slug           string            `json:"slug"`
	Content        string            `json:"content"`
	MainImageURI   string            `json:"mainImageURI"`
	AuthorId       int               `json:"authorId"`
	AuthorDetail   *AuthorResponse   `gorm:"embedded;embeddedPrefix:AuthorDetail_" json:"authorDetail,omitempty"`
	CategoryDetail *CategoryResponse `gorm:"embedded;embeddedPrefix:CategoryDetail_" json:"categoryDetail,omitempty"`
	LikeCount      int64             `json:"likeCount"`
	Status         int               `json:"status"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}

type PostUploadResponse struct {
	Url         string `json:"url"`
	IsTemporary int16  `json:"isTemporary"`
}
