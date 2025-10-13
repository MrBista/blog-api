package dto

type CommentRequest struct {
	PostId   int    `json:"postId"`
	ParentId int    `json:"parentId"`
	Content  string `json:"content" validate:"required"`
}
