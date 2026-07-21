package models

import (
	"time"

	"gorm.io/datatypes"
)

type SourceSyncLog struct {
	Common
	LibraryID    string     `gorm:"column:library_id;type:text;index;not null" json:"libraryId"`
	SourceType   string     `gorm:"column:source_type;type:text;not null" json:"sourceType"`
	Status       string     `gorm:"column:status;type:text;not null" json:"status"`
	ItemsCreated int        `gorm:"column:items_created;default:0" json:"itemsCreated"`
	ItemsUpdated int        `gorm:"column:items_updated;default:0" json:"itemsUpdated"`
	ItemsRemoved int        `gorm:"column:items_removed;default:0" json:"itemsRemoved"`
	StartedAt    time.Time  `gorm:"column:started_at;not null" json:"startedAt"`
	CompletedAt  *time.Time `gorm:"column:completed_at" json:"completedAt,omitempty"`
	ErrorMessage *string    `gorm:"column:error_message;type:text" json:"errorMessage,omitempty"`
}

func (SourceSyncLog) TableName() string { return "source_sync_logs" }

type SourceConfig struct {
	Common
	SourceType string         `gorm:"column:source_type;type:text;not null;uniqueIndex" json:"sourceType"`
	Enabled    bool           `gorm:"column:enabled;default:true" json:"enabled"`
	Config     datatypes.JSON `gorm:"column:config;type:jsonb;not null" json:"config"`
	LastSyncAt *time.Time     `gorm:"column:last_sync_at" json:"lastSyncAt,omitempty"`
}

func (SourceConfig) TableName() string { return "source_configs" }
