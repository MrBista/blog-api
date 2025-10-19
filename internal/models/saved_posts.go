package models

import "time"

type SavedPost struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64      `gorm:"not null;index:idx_user_saved_posts" json:"userId"`
	PostID        int64      `gorm:"not null;index:idx_post_saved" json:"postId"`
	ReadingListID int64      `gorm:"not null;index:idx_reading_list_posts" json:"readingListId"`
	Notes         *string    `gorm:"type:text" json:"notes"`
	IsRead        bool       `gorm:"default:false;index:idx_is_read" json:"isRead"`
	ReadAt        *time.Time `json:"readAt"`
	CreatedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`

	// Relations
	// User        User        `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	// Post        Post        `gorm:"foreignKey:PostID;references:ID" json:"post,omitempty"`
	// ReadingList ReadingList `gorm:"foreignKey:ReadingListID;references:ID" json:"readingList,omitempty"`
}

func (s *SavedPost) TableName() string {
	return "saved_posts"
}
