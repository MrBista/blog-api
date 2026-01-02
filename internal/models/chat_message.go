package models

import "time"

type ChatMessage struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	RoomId    int64     `gorm:"column:room_id"`
	SenderId  int64     `gorm:"column:sender_id"`
	Message   string    `gorm:"column:message"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (c *ChatMessage) TableName() string {
	return "chat_messages"
}
