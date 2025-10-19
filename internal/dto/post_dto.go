package dto

import "time"

type CreatePostRequest struct {
	Id         int    `json:"id,omitempty"`
	Title      string `json:"title" validate:"required"`
	Content    string `json:"content" validate:"required"`
	CategoryId int    `json:"categoryId" validate:"required"`
	ImgUrl     string `json:"imgUrl"`
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

// ReadingListDTO untuk response list reading list
type ReadingListDTO struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userId"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	IsDefault   bool      `json:"isDefault"`
	Color       *string   `json:"color"`
	Icon        *string   `json:"icon"`
	OrderIndex  int       `json:"orderIndex"`
	TotalPosts  int       `json:"totalPosts"`
	UnreadCount int       `json:"unreadCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateReadingListRequest untuk membuat reading list baru
type CreateReadingListRequest struct {
	Name        string  `json:"name" validate:"required,max=100"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
	Icon        *string `json:"icon"`
	OrderIndex  int     `json:"orderIndex"`
}

// UpdateReadingListRequest untuk update reading list
type UpdateReadingListRequest struct {
	Name        *string `json:"name" validate:"omitempty,max=100"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
	Icon        *string `json:"icon"`
	OrderIndex  *int    `json:"orderIndex"`
}

// CreateSavedPostRequest untuk menyimpan post ke reading list
type CreateSavedPostRequest struct {
	PostID        int64   `json:"postId" validate:"required"`
	ReadingListID int64   `json:"readingListId" validate:"required"`
	Notes         *string `json:"notes"`
}

// SavedPostDTO untuk response saved post
type SavedPostDTO struct {
	ID            int64      `json:"id"`
	UserID        int64      `json:"userId"`
	PostID        int64      `json:"postId"`
	ReadingListID int64      `json:"readingListId"`
	Notes         *string    `json:"notes"`
	IsRead        bool       `json:"isRead"`
	ReadAt        *time.Time `json:"readAt"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`

	// Embedded post info (optional)
	Post *SavedPostInfo `json:"post,omitempty"`
}

// SavedPostInfo info singkat post yang disimpan
type SavedPostInfo struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	Slug         string  `json:"slug"`
	MainImageURI *string `json:"mainImageUri"`
	AuthorName   string  `json:"authorName"`
	CategoryName *string `json:"categoryName"`
}

// UpdateSavedPostRequest untuk update saved post (misal mark as read, update notes)
type UpdateSavedPostRequest struct {
	Notes  *string `json:"notes"`
	IsRead *bool   `json:"isRead"`
}
