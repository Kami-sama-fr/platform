package models

import (

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Category struct {
	Common
	Name        string         `gorm:"column:name;type:text;not null" json:"name"`
	Slug        string         `gorm:"column:slug;type:text;uniqueIndex;not null" json:"slug"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	ParentID    *string        `gorm:"column:parent_id;type:text;index" json:"parentId,omitempty"`
	SortOrder   int            `gorm:"column:sort_order;default:0" json:"sortOrder"`
	IsActive    bool           `gorm:"column:is_active;default:true" json:"isActive"`
	Metadata    datatypes.JSON `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (Category) TableName() string { return "categories" }
