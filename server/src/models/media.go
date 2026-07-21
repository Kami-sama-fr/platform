package models

import (
	"time"

	"gorm.io/datatypes"
)

type MediaAsset struct {
	Common
	AnimeId     string         `gorm:"column:anime_id;type:text;index" json:"animeId"`
	EpisodeId   *string        `gorm:"column:episode_id;type:text;index" json:"episodeId,omitempty"`
	Type        string         `gorm:"column:type;type:text;not null" json:"type"`
	Filename    string         `gorm:"column:filename;type:text;not null" json:"filename"`
	OriginalUrl string         `gorm:"column:original_url;type:text" json:"originalUrl"`
	MimeType    string         `gorm:"column:mime_type;type:text" json:"mimeType"`
	Size        int64          `gorm:"column:size;not null" json:"size"`
	Duration    float64        `gorm:"column:duration;type:real" json:"duration"`
	Width       int            `gorm:"column:width;default:0" json:"width"`
	Height      int            `gorm:"column:height;default:0" json:"height"`
	Language    string         `gorm:"column:language;type:text" json:"language"`
	Quality     string         `gorm:"column:quality;type:text" json:"quality"`
	Codec       string         `gorm:"column:codec;type:text" json:"codec"`
	CdnUrl      string         `gorm:"column:cdn_url;type:text" json:"cdnUrl"`
	Status      string         `gorm:"column:status;type:text;not null;default:'pending'" json:"status"`
	Metadata    datatypes.JSON `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
}

func (MediaAsset) TableName() string { return "media_assets" }

type EncodingJob struct {
	Common
	MediaAssetId string         `gorm:"column:media_asset_id;type:text;index;not null" json:"mediaAssetId"`
	Status       string         `gorm:"column:status;type:text;not null;default:'pending'" json:"status"`
	Profile      string         `gorm:"column:profile;type:text" json:"profile"`
	Progress     int            `gorm:"column:progress;default:0" json:"progress"`
	OutputUrl    string         `gorm:"column:output_url;type:text" json:"outputUrl"`
	StartedAt    *time.Time     `gorm:"column:started_at" json:"startedAt,omitempty"`
	CompletedAt  *time.Time     `gorm:"column:completed_at" json:"completedAt,omitempty"`
	FailedAt     *time.Time     `gorm:"column:failed_at" json:"failedAt,omitempty"`
	ErrorMessage *string        `gorm:"column:error_message;type:text" json:"errorMessage,omitempty"`
	Metadata     datatypes.JSON `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
}

func (EncodingJob) TableName() string { return "encoding_jobs" }

type CdnAsset struct {
	Common
	OriginalUrl string `gorm:"column:original_url;type:text;not null" json:"originalUrl"`
	CdnUrl      string `gorm:"column:cdn_url;type:text;not null" json:"cdnUrl"`
	FileType    string `gorm:"column:file_type;type:text" json:"fileType"`
	FileSize    int64  `gorm:"column:file_size;not null" json:"fileSize"`
	Status      string `gorm:"column:status;type:text;not null;default:'active'" json:"status"`
}

func (CdnAsset) TableName() string { return "cdn_assets" }
