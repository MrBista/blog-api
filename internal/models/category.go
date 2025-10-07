package models

type Category struct {
	ID          int64   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string  `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Slug        string  `gorm:"column:slug;type:varchar(100);unique;not null" json:"slug"`
	Description *string `gorm:"column:description;type:text" json:"description,omitempty"`
	ParentID    *int64  `gorm:"column:parent_id" json:"parentId,omitempty"`

	// Relations
	// Parent   *Category  `gorm:"foreignKey:ParentID;references:ID" json:"parent,omitempty"`
	// Children []Category `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"children,omitempty"`
	// Posts    []Post     `gorm:"foreignKey:CategoryID;references:ID" json:"posts,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}
