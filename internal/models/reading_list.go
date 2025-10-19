package models

import "time"

type ReadingList struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64     `gorm:"not null;index:idx_user_reading_lists" json:"userId"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description"`
	IsDefault   bool      `gorm:"default:false" json:"isDefault"`
	Color       *string   `gorm:"type:varchar(20)" json:"color"`
	Icon        *string   `gorm:"type:varchar(50)" json:"icon"`
	OrderIndex  int       `gorm:"default:0" json:"orderIndex"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`

	// Relations
	// User       User        `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	// SavedPosts []SavedPost `gorm:"foreignKey:ReadingListID" json:"savedPosts,omitempty"`
}

func (r *ReadingList) TableName() string {
	return "reading_lists"
}
