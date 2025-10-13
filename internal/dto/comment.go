package dto

type CommentRequest struct {
	PostId   int    `json:"postId" validate:"required; numeric"`
	ParentId int    `json:"parentId"`
	Content  string `json:"content" validate:"required"`
}
