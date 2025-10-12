package dto

type PostFilterRequest struct {
	Title      string `json:"title" query:"title"`
	CategoryID int    `json:"category_id" query:"category_id"`
	AuthorID   int    `json:"author_id" query:"author_id"`
	Status     int    `json:"status" query:"status"`
	PaginationParams
}

type UserFilterRequest struct {
	Email    string `json:"email" query:"email"`
	Username string `json:"username" query:"username"`
	Role     int    `json:"role" query:"role"`
	PaginationParams
}

type CategoryFilterRequest struct {
	Name string `json:"name" query:"name"`
	PaginationParams
}
