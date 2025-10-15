package models

import "time"

type PostAsset struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID      int64     `gorm:"not null;index" json:"postId"`
	AssetURI    string    `gorm:"type:varchar(500);not null" json:"assetUri"`
	Type        uint8     `gorm:"type:tinyint;default:1;comment:'1=image, 2=video, 3=file'" json:"type"`
	Caption     *string   `gorm:"type:varchar(255)" json:"caption,omitempty"`
	OrderIndex  int       `gorm:"default:0" json:"orderIndex"`
	IsTemporary uint8     `gorm:"type:tinyint;default:1;comment:'1=true, 0=false'" json:"isTemporary"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`

	// Relationship (optional)
	// Post *Post `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"post,omitempty"`
}

func (PostAsset) TableName() string {
	return "post_assets"
}
