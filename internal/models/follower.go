package models

import "time"

type Follower struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	FollowerID  uint64    `gorm:"column:follower_id;not null;index:idx_follower" json:"followerId"`    // User yang melakukan follow
	FollowingID uint64    `gorm:"column:following_id;not null;index:idx_following" json:"followingId"` // User yang di-follow
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
}

func (f *Follower) TableName() string {
	return "followers"
}
