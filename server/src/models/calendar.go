package models

import (
	"time"

	"gorm.io/datatypes"
)

type CalendarEvent struct {
	Common
	Title       string         `gorm:"column:title;type:text;not null" json:"title"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	EventType   string         `gorm:"column:event_type;type:text;not null;default:'release'" json:"eventType"`
	AnimeID     *string        `gorm:"column:anime_id;type:text;index" json:"animeId,omitempty"`
	EpisodeID   *string        `gorm:"column:episode_id;type:text" json:"episodeId,omitempty"`
	StartAt     time.Time      `gorm:"column:start_at;not null" json:"startAt"`
	EndAt       *time.Time     `gorm:"column:end_at" json:"endAt,omitempty"`
	AllDay      bool           `gorm:"column:all_day;default:false" json:"allDay"`
	Status      string         `gorm:"column:status;type:text;not null;default:'scheduled'" json:"status"`
	Color       string         `gorm:"column:color;type:text" json:"color"`
	Metadata    datatypes.JSON `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
}

func (CalendarEvent) TableName() string { return "calendar_events" }

type Premiere struct {
	Common
	AnimeID    string     `gorm:"column:anime_id;type:text;index;not null" json:"animeId"`
	Season     int        `gorm:"column:season;default:1" json:"season"`
	EpisodeNum int        `gorm:"column:episode_num;default:1" json:"episodeNum"`
	ScheduledAt time.Time `gorm:"column:scheduled_at;not null" json:"scheduledAt"`
	Region     string     `gorm:"column:region;type:text;not null;default:'global'" json:"region"`
	Status     string     `gorm:"column:status;type:text;not null;default:'scheduled'" json:"status"`
}

func (Premiere) TableName() string { return "premieres" }
