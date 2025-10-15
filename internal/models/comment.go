package models

import "time"

type Comment struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PostID    int64     `gorm:"column:post_id;not null" json:"postId"`
	UserID    *int64    `gorm:"column:user_id" json:"userId,omitempty"`
	ParentID  *int64    `gorm:"column:parent_id" json:"parentId"`
	Name      *string   `gorm:"column:name;type:varchar(150)" json:"name,omitempty"`
	Email     *string   `gorm:"column:email;type:varchar(150)" json:"email,omitempty"`
	Content   string    `gorm:"column:content;type:text;not null" json:"content"`
	Status    int8      `gorm:"column:status;type:tinyint;default:1;comment:0=inactive,1=active,2=deleted" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (c *Comment) TableName() string {
	return "comments"
}
