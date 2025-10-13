package models

import "time"

type Like struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"column:user_id;not null" json:"userId"`
	TargetType int8      `gorm:"column:target_type;type:tinyint;not null;default:1;comment:1=posts,2=comments" json:"targetType"`
	TargetID   int64     `gorm:"column:target_id;not null" json:"targetId"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (l *Like) TableName() string {
	return "likes"
}
