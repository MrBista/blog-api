package dto

import "time"

type CommentWithUserResponse struct {
	ID        int64              `json:"id"`
	PostID    int64              `json:"postId"`
	UserID    *int64             `json:"userId,omitempty"`
	ParentID  *int64             `json:"parentId,omitempty"`
	Name      *string            `json:"name,omitempty"`
	Email     *string            `json:"email,omitempty"`
	Content   string             `json:"content"`
	Status    int8               `json:"status"`
	CreatedAt time.Time          `json:"createdAt"`
	User      *UserBriefResponse `json:"user,omitempty"`
}

type UserBriefResponse struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	Username        *string `json:"username,omitempty"`
	Email           string  `json:"email"`
	ProfileImageURI *string `json:"profileImageUri"`
}
