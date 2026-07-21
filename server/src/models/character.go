package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Character struct {
	Common
	Name        string         `gorm:"column:name;type:text;not null" json:"name"`
	Slug        string         `gorm:"column:slug;type:text;uniqueIndex;not null" json:"slug"`
	ImageUrl    string         `gorm:"column:image_url;type:text" json:"imageUrl"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	Gender      string         `gorm:"column:gender;type:text" json:"gender"`
	Metadata    datatypes.JSON `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (Character) TableName() string { return "characters" }
