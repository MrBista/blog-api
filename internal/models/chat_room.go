package models

import "time"

type ChatRoom struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	ArticleId int64     `gorm:"column:post_id"`
	ReaderId  int64     `gorm:"column:reader_id"`
	WriterId  int64     `gorm:"columen:writer_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (c *ChatRoom) TableName() string {
	return "chat_rooms"
}
