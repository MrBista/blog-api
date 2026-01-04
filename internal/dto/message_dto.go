package dto

type ChatRoomDto struct {
	PostId int `json:"post_id" validate:"required"`
}
