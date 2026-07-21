package models

import "gorm.io/gorm"

type Genre struct {
	Common
	Name        string         `gorm:"column:name;type:text;uniqueIndex;not null" json:"name"`
	Slug        string         `gorm:"column:slug;type:text;uniqueIndex;not null" json:"slug"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (Genre) TableName() string { return "genres" }
