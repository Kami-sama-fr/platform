package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Episode struct {
	Common
	AnimeID      string         `gorm:"column:anime_id;type:text;index;not null" json:"animeId"`
	SeasonID     *string        `gorm:"column:season_id;type:text;index" json:"seasonId,omitempty"`
	Number       int            `gorm:"column:number;not null" json:"number"`
	Title        string         `gorm:"column:title;type:text;not null" json:"title"`
	Synopsis     string         `gorm:"column:synopsis;type:text" json:"synopsis"`
	ThumbnailUrl string         `gorm:"column:thumbnail_url;type:text" json:"thumbnailUrl"`
	Duration     float64        `gorm:"column:duration;type:real" json:"duration"`
	AirDate      *time.Time     `gorm:"column:air_date" json:"airDate,omitempty"`
	IsSubbed     bool           `gorm:"column:is_subbed;default:true" json:"isSubbed"`
	IsDubbed     bool           `gorm:"column:is_dubbed;default:false" json:"isDubbed"`
	Metadata     datatypes.JSON `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (Episode) TableName() string { return "episodes" }
