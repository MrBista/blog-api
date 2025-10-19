package models

import "time"

type User struct {
	ID              int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name            string     `gorm:"column:name;type:varchar(150);not null" json:"name"`
	Username        string     `gorm:"column:username;type:varchar(100);unique;not null" json:"username"`
	Email           string     `gorm:"column:email;type:varchar(150);unique;not null" json:"email"`
	Password        string     `gorm:"column:password;type:varchar(255)" json:"password"`
	Bio             *string    `gorm:"column:bio;type:text" json:"bio,omitempty"`
	ProfileImageURI *string    `gorm:"column:profile_image_uri;type:varchar(500)" json:"profileImageUri,omitempty"`
	Role            int        `gorm:"column:role;default:0;comment:0=reader,1=editor,2=author,3=admin" json:"role"`
	IsSubscribed    bool       `gorm:"column:is_subscribed;default:false" json:"isSubscribed"`
	SubscriptionEnd *time.Time `gorm:"column:subscription_end" json:"subscriptionEnd"`
	Status          int        `gorm:"column:status;default:0;comment:0=inactive,1=active,2=archived,3=banned" json:"status"`
	AuthProvider    string     `gorm:"default:'local'"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// Relations
	// Posts []Post `gorm:"foreignKey:AuthorID;references:ID" json:"posts,omitempty"`
}

func (User) TableName() string {
	return "users"
}
