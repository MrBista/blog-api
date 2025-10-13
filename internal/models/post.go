package models

import "time"

type Post struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title          string     `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Slug           string     `gorm:"column:slug;type:varchar(255);unique;not null" json:"slug"`
	Content        string     `gorm:"column:content;type:longtext;not null" json:"content"`
	MainImageURI   *string    `gorm:"column:main_image_uri;type:varchar(500)" json:"mainImageUri,omitempty"`
	AuthorID       int64      `gorm:"column:author_id;not null" json:"authorId"`
	CategoryID     *int64     `gorm:"column:category_id" json:"categoryId,omitempty"`
	Status         uint8      `gorm:"column:status;default:0;comment:0='inactive',1='draft',2='review',3='published',4='archived'" json:"status"`
	IsFeatured     bool       `gorm:"column:is_featured;default:false" json:"isFeatured"`
	ViewCount      int        `gorm:"column:view_count;default:0" json:"viewCount"`
	SeoTitle       *string    `gorm:"column:seo_title;type:varchar(255)" json:"seoTitle,omitempty"`
	SeoDescription *string    `gorm:"column:seo_description;type:varchar(255)" json:"seoDescription,omitempty"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	PublishedAt    *time.Time `gorm:"column:published_at" json:"publishedAt,omitempty"`

	// Relations
	Author   *User     `gorm:"foreignKey:AuthorID;references:ID" json:"author,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
	Likes    []Like    `gorm:"foreignKey:TargetID;references:ID" json:"likes,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID;references:ID" json:"comments,omitempty"`

	// LikeCount int64 `gorm:"-" json:"likeCount"`
	LikeCount int64 `gorm:"->;column:like_count" json:"likeCount"`
}

func (p *Post) TableName() string {
	return "posts"
}
